package utils

import (
	"context"
	"crypto/rand"
	"encoding/base32"
	"errors"
	"fmt"
	"html/template"
	"regexp"
	"strconv"
	"strings"
	"time"

	"wms-service/internal/domain"
	customError "wms-service/pkg/error"

	"github.com/gin-gonic/gin"
	validatorPkg "github.com/go-playground/validator/v10"
	"github.com/oklog/ulid"
	"github.com/omniful/go_commons/config"
	"github.com/omniful/go_commons/constants"
	"github.com/omniful/go_commons/email"
	"github.com/omniful/go_commons/env"
	oerror "github.com/omniful/go_commons/error"
	"github.com/omniful/go_commons/jwt/private"
	"github.com/omniful/go_commons/log"
	cons "github.com/omniful/tenant-service/constants"
	"github.com/xlzd/gotp"
)

func GetDefaultHeaders(ctx context.Context) (headers map[string][]string) {
	headers = make(map[string][]string, 0)
	val := ctx.Value(constants.JWTHeader).(string)
	headers[constants.JWTHeader] = []string{val}

	return
}

func GenerateULID() string {
	entropy := ulid.Monotonic(rand.Reader, 0)
	id := ulid.MustNew(ulid.Timestamp(time.Now()), entropy)
	return id.String()
}

func GetValidationErrors(ctx context.Context, validationErr error) (cusErr oerror.CustomError) {
	var ve validatorPkg.ValidationErrors
	if errors.As(validationErr, &ve) {
		errorsMap := make([]ErrorMapStruct, 0)
		for _, errs := range ve {
			errorsMap = append(errorsMap, ErrorMapStruct{Field: errs.Field(), Tag: errs.Tag()})
		}

		for _, err := range errorsMap {
			cusErr = customError.InvalidRequest(ctx, err.Field+" "+err.Tag)
			return
		}
	}

	return
}

func Remove(slice []string, element string) []string {
	for i, v := range slice {
		if v == element {
			// Remove the element by slicing the slice around it
			return append(slice[:i], slice[i+1:]...)
		}
	}
	return slice
}

type ErrorMapStruct struct {
	Field string `json:"field"`
	Tag   string `json:"tag"`
}

func ExtractFirstCaptureGroup(regex *regexp.Regexp, input string) (capturedValue string) {
	match := regex.FindStringSubmatch(input)
	if len(match) < 2 {
		return "" // or return an error if a match is not found
	}
	return match[1]
}

func GetUserDetails(ctx context.Context) (userDetails domain.UserDetails, cusErr oerror.CustomError) {
	userID, cusErr := private.GetUserID(ctx)
	if cusErr.Exists() {
		return
	}
	userDetails.UserID = userID

	tenantID, cusErr := private.GetTenantID(ctx)
	if cusErr.Exists() {
		return
	}
	userDetails.TenantID = tenantID

	userName, cusErr := private.GetUserName(ctx)
	if cusErr.Exists() {
		return
	}
	userDetails.UserName = userName

	tenantName, cusErr := private.GetTenantName(ctx)
	if cusErr.Exists() {
		return
	}
	userDetails.TenantName = tenantName

	userEmail, cusErr := private.GetUserEmail(ctx)
	if cusErr.Exists() {
		return
	}
	userDetails.UserEmail = userEmail

	return
}

func GetNameSpace(ctx context.Context) string {
	return config.GetString(ctx, "service.name")
}

func SendEmail(ctx context.Context, tenantID string, client email.EmailClientV2, subject, attachmentUrl, message, userEmail, htmlFile string, ccEmails ...string) (cusErr oerror.CustomError) {
	messageData := email.Message{
		Subject:      subject,
		Template:     nil,
		TemplateData: nil,
	}

	recipient := email.Recipient{
		ToEmails: []string{userEmail},
	}
	recipient.CcEmails = append(recipient.CcEmails, ccEmails...)

	messageData.TemplateData = struct {
		UserEmail     string
		AttachmentURL string
		Message       string
		CreatedAt     string
	}{
		UserEmail:     userEmail,
		AttachmentURL: attachmentUrl,
		Message:       message,
		CreatedAt:     time.Now().Format(cons.DefaultTimeFormat),
	}

	t, err := template.ParseFiles(htmlFile)
	if err != nil {
		log.Errorf("error in sending email to sqs :: %v", err)
		cusErr = oerror.NewCustomError(customError.ParseFilesError, fmt.Sprintf("parse file error"))
		return
	}

	messageData.Template = t
	err = client.SendEmail(ctx, tenantID, messageData, recipient)
	if err != nil {
		log.Errorf("Error sending email: %v", err)
		cusErr = oerror.NewCustomError(customError.BadRequest, err.Error())
		return
	}

	return
}

func ParsePhoneNumber(phoneNumber string) (countryCode, countryCallingCode, mobileNumber string) {
	splitNumber := strings.Split(phoneNumber, "-")

	if len(splitNumber) == 3 {
		return splitNumber[0], splitNumber[1], splitNumber[2]
	}

	return "", "", phoneNumber
}

func PhoneNumber(countryCode, countryCallingCode, mobileNumber string) (phoneNumber string) {
	return countryCode + "-" + countryCallingCode + "-" + mobileNumber
}

func CapitalizeFirstChar(input string) string {
	if len(input) == 0 {
		return input
	}

	return strings.ToUpper(string(input[0])) + input[1:]
}

func GenerateOTP(ctx context.Context) (OTP string, cusErr oerror.CustomError) {
	baseSalt := config.GetString(ctx, "otp.salt")
	currentTime := time.Now()
	dynamicSalt := fmt.Sprintf("%s%d", baseSalt, currentTime.UnixNano())
	encodedSalt := base32.StdEncoding.EncodeToString([]byte(dynamicSalt))
	encodedSalt = strings.TrimRight(encodedSalt, "=")

	otp := gotp.NewDefaultTOTP(encodedSalt)
	otp.At(time.Now().UnixNano())
	OTP = otp.Now()
	return
}

func ExtractTenantIDFromParamsAndContext(ctx *gin.Context) (uint64, oerror.CustomError) {
	tenantIDStr, cusErr := private.GetTenantID(ctx)
	if cusErr.Exists() {
		return 0, cusErr
	}

	tenantIDParam := ctx.Param(cons.TenantID)
	if tenantIDStr != tenantIDParam {
		cusErr = oerror.NewCustomError(
			customError.BadRequest,
			fmt.Sprintf("tenantID in path param and header are not the same :: [%s] != [%s]", tenantIDParam, tenantIDStr),
		)
		return 0, cusErr
	}

	tenantID, cusErr := ParseStringToUint64(ctx, tenantIDStr)
	if cusErr.Exists() {
		return 0, cusErr
	}

	return tenantID, oerror.CustomError{}
}

func ParseParamToUint64(ctx *gin.Context, param string) (uint64, oerror.CustomError) {
	return ParseStringToUint64(ctx, ctx.Param(param))
}

func ParseStringToUint64(ctx context.Context, str string) (uint64, oerror.CustomError) {
	logTag := fmt.Sprintf("RequestID: %s Function: Process (ParseStringToUint64)", env.GetRequestID(ctx))

	num, err := strconv.ParseUint(str, 10, 64)
	if err != nil {
		return 0, oerror.NewCustomError(
			customError.ParseIntError,
			fmt.Sprintf("%s unable to parse str :: [%s] | err :: [%s]", logTag, str, err.Error()),
		)
	}

	return num, oerror.CustomError{}
}

func TrimLowerArr(input []string) []string {
	for i, v := range input {
		input[i] = TrimLower(v)
	}
	return input
}

func TrimLower(input string) string {
	return strings.ToLower(strings.TrimSpace(input))
}

func GetLogTag(ctx context.Context, funcName string) string {
	return fmt.Sprintf("Request ID: %s Function: %s, ", env.GetRequestID(ctx), funcName)
}
