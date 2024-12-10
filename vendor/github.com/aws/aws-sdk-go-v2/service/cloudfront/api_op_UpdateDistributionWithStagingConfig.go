// Code generated by smithy-go-codegen DO NOT EDIT.

package cloudfront

import (
	"context"
	"fmt"
	awsmiddleware "github.com/aws/aws-sdk-go-v2/aws/middleware"
	"github.com/aws/aws-sdk-go-v2/service/cloudfront/types"
	"github.com/aws/smithy-go/middleware"
	smithyhttp "github.com/aws/smithy-go/transport/http"
)

// Copies the staging distribution's configuration to its corresponding primary
// distribution. The primary distribution retains its Aliases (also known as
// alternate domain names or CNAMEs) and ContinuousDeploymentPolicyId value, but
// otherwise its configuration is overwritten to match the staging distribution.
//
// You can use this operation in a continuous deployment workflow after you have
// tested configuration changes on the staging distribution. After using a
// continuous deployment policy to move a portion of your domain name's traffic to
// the staging distribution and verifying that it works as intended, you can use
// this operation to copy the staging distribution's configuration to the primary
// distribution. This action will disable the continuous deployment policy and move
// your domain's traffic back to the primary distribution.
//
// This API operation requires the following IAM permissions:
//
// [GetDistribution]
//
// [UpdateDistribution]
//
// [GetDistribution]: https://docs.aws.amazon.com/cloudfront/latest/APIReference/API_GetDistribution.html
// [UpdateDistribution]: https://docs.aws.amazon.com/cloudfront/latest/APIReference/API_UpdateDistribution.html
func (c *Client) UpdateDistributionWithStagingConfig(ctx context.Context, params *UpdateDistributionWithStagingConfigInput, optFns ...func(*Options)) (*UpdateDistributionWithStagingConfigOutput, error) {
	if params == nil {
		params = &UpdateDistributionWithStagingConfigInput{}
	}

	result, metadata, err := c.invokeOperation(ctx, "UpdateDistributionWithStagingConfig", params, optFns, c.addOperationUpdateDistributionWithStagingConfigMiddlewares)
	if err != nil {
		return nil, err
	}

	out := result.(*UpdateDistributionWithStagingConfigOutput)
	out.ResultMetadata = metadata
	return out, nil
}

type UpdateDistributionWithStagingConfigInput struct {

	// The identifier of the primary distribution to which you are copying a staging
	// distribution's configuration.
	//
	// This member is required.
	Id *string

	// The current versions ( ETag values) of both primary and staging distributions.
	// Provide these in the following format:
	//
	//     <primary ETag>, <staging ETag>
	IfMatch *string

	// The identifier of the staging distribution whose configuration you are copying
	// to the primary distribution.
	StagingDistributionId *string

	noSmithyDocumentSerde
}

type UpdateDistributionWithStagingConfigOutput struct {

	// A distribution tells CloudFront where you want content to be delivered from,
	// and the details about how to track and manage content delivery.
	Distribution *types.Distribution

	// The current version of the primary distribution (after it's updated).
	ETag *string

	// Metadata pertaining to the operation's result.
	ResultMetadata middleware.Metadata

	noSmithyDocumentSerde
}

func (c *Client) addOperationUpdateDistributionWithStagingConfigMiddlewares(stack *middleware.Stack, options Options) (err error) {
	if err := stack.Serialize.Add(&setOperationInputMiddleware{}, middleware.After); err != nil {
		return err
	}
	err = stack.Serialize.Add(&awsRestxml_serializeOpUpdateDistributionWithStagingConfig{}, middleware.After)
	if err != nil {
		return err
	}
	err = stack.Deserialize.Add(&awsRestxml_deserializeOpUpdateDistributionWithStagingConfig{}, middleware.After)
	if err != nil {
		return err
	}
	if err := addProtocolFinalizerMiddlewares(stack, options, "UpdateDistributionWithStagingConfig"); err != nil {
		return fmt.Errorf("add protocol finalizers: %v", err)
	}

	if err = addlegacyEndpointContextSetter(stack, options); err != nil {
		return err
	}
	if err = addSetLoggerMiddleware(stack, options); err != nil {
		return err
	}
	if err = addClientRequestID(stack); err != nil {
		return err
	}
	if err = addComputeContentLength(stack); err != nil {
		return err
	}
	if err = addResolveEndpointMiddleware(stack, options); err != nil {
		return err
	}
	if err = addComputePayloadSHA256(stack); err != nil {
		return err
	}
	if err = addRetry(stack, options); err != nil {
		return err
	}
	if err = addRawResponseToMetadata(stack); err != nil {
		return err
	}
	if err = addRecordResponseTiming(stack); err != nil {
		return err
	}
	if err = addSpanRetryLoop(stack, options); err != nil {
		return err
	}
	if err = addClientUserAgent(stack, options); err != nil {
		return err
	}
	if err = smithyhttp.AddErrorCloseResponseBodyMiddleware(stack); err != nil {
		return err
	}
	if err = smithyhttp.AddCloseResponseBodyMiddleware(stack); err != nil {
		return err
	}
	if err = addSetLegacyContextSigningOptionsMiddleware(stack); err != nil {
		return err
	}
	if err = addTimeOffsetBuild(stack, c); err != nil {
		return err
	}
	if err = addUserAgentRetryMode(stack, options); err != nil {
		return err
	}
	if err = addOpUpdateDistributionWithStagingConfigValidationMiddleware(stack); err != nil {
		return err
	}
	if err = stack.Initialize.Add(newServiceMetadataMiddleware_opUpdateDistributionWithStagingConfig(options.Region), middleware.Before); err != nil {
		return err
	}
	if err = addRecursionDetection(stack); err != nil {
		return err
	}
	if err = addRequestIDRetrieverMiddleware(stack); err != nil {
		return err
	}
	if err = addResponseErrorMiddleware(stack); err != nil {
		return err
	}
	if err = addRequestResponseLogging(stack, options); err != nil {
		return err
	}
	if err = addDisableHTTPSMiddleware(stack, options); err != nil {
		return err
	}
	if err = addSpanInitializeStart(stack); err != nil {
		return err
	}
	if err = addSpanInitializeEnd(stack); err != nil {
		return err
	}
	if err = addSpanBuildRequestStart(stack); err != nil {
		return err
	}
	if err = addSpanBuildRequestEnd(stack); err != nil {
		return err
	}
	return nil
}

func newServiceMetadataMiddleware_opUpdateDistributionWithStagingConfig(region string) *awsmiddleware.RegisterServiceMetadata {
	return &awsmiddleware.RegisterServiceMetadata{
		Region:        region,
		ServiceID:     ServiceID,
		OperationName: "UpdateDistributionWithStagingConfig",
	}
}
