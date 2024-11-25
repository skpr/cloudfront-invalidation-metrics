// Code generated by smithy-go-codegen DO NOT EDIT.

package cloudwatch

import (
	"context"
	"fmt"
	awsmiddleware "github.com/aws/aws-sdk-go-v2/aws/middleware"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatch/types"
	"github.com/aws/smithy-go/middleware"
	smithyhttp "github.com/aws/smithy-go/transport/http"
)

// Creates a Contributor Insights rule. Rules evaluate log events in a CloudWatch
// Logs log group, enabling you to find contributor data for the log events in that
// log group. For more information, see [Using Contributor Insights to Analyze High-Cardinality Data].
//
// If you create a rule, delete it, and then re-create it with the same name,
// historical data from the first time the rule was created might not be available.
//
// [Using Contributor Insights to Analyze High-Cardinality Data]: https://docs.aws.amazon.com/AmazonCloudWatch/latest/monitoring/ContributorInsights.html
func (c *Client) PutInsightRule(ctx context.Context, params *PutInsightRuleInput, optFns ...func(*Options)) (*PutInsightRuleOutput, error) {
	if params == nil {
		params = &PutInsightRuleInput{}
	}

	result, metadata, err := c.invokeOperation(ctx, "PutInsightRule", params, optFns, c.addOperationPutInsightRuleMiddlewares)
	if err != nil {
		return nil, err
	}

	out := result.(*PutInsightRuleOutput)
	out.ResultMetadata = metadata
	return out, nil
}

type PutInsightRuleInput struct {

	// The definition of the rule, as a JSON object. For details on the valid syntax,
	// see [Contributor Insights Rule Syntax].
	//
	// [Contributor Insights Rule Syntax]: https://docs.aws.amazon.com/AmazonCloudWatch/latest/monitoring/ContributorInsights-RuleSyntax.html
	//
	// This member is required.
	RuleDefinition *string

	// A unique name for the rule.
	//
	// This member is required.
	RuleName *string

	// The state of the rule. Valid values are ENABLED and DISABLED.
	RuleState *string

	// A list of key-value pairs to associate with the Contributor Insights rule. You
	// can associate as many as 50 tags with a rule.
	//
	// Tags can help you organize and categorize your resources. You can also use them
	// to scope user permissions, by granting a user permission to access or change
	// only the resources that have certain tag values.
	//
	// To be able to associate tags with a rule, you must have the
	// cloudwatch:TagResource permission in addition to the cloudwatch:PutInsightRule
	// permission.
	//
	// If you are using this operation to update an existing Contributor Insights
	// rule, any tags you specify in this parameter are ignored. To change the tags of
	// an existing rule, use [TagResource].
	//
	// [TagResource]: https://docs.aws.amazon.com/AmazonCloudWatch/latest/APIReference/API_TagResource.html
	Tags []types.Tag

	noSmithyDocumentSerde
}

type PutInsightRuleOutput struct {
	// Metadata pertaining to the operation's result.
	ResultMetadata middleware.Metadata

	noSmithyDocumentSerde
}

func (c *Client) addOperationPutInsightRuleMiddlewares(stack *middleware.Stack, options Options) (err error) {
	if err := stack.Serialize.Add(&setOperationInputMiddleware{}, middleware.After); err != nil {
		return err
	}
	err = stack.Serialize.Add(&awsAwsquery_serializeOpPutInsightRule{}, middleware.After)
	if err != nil {
		return err
	}
	err = stack.Deserialize.Add(&awsAwsquery_deserializeOpPutInsightRule{}, middleware.After)
	if err != nil {
		return err
	}
	if err := addProtocolFinalizerMiddlewares(stack, options, "PutInsightRule"); err != nil {
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
	if err = addOpPutInsightRuleValidationMiddleware(stack); err != nil {
		return err
	}
	if err = stack.Initialize.Add(newServiceMetadataMiddleware_opPutInsightRule(options.Region), middleware.Before); err != nil {
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

func newServiceMetadataMiddleware_opPutInsightRule(region string) *awsmiddleware.RegisterServiceMetadata {
	return &awsmiddleware.RegisterServiceMetadata{
		Region:        region,
		ServiceID:     ServiceID,
		OperationName: "PutInsightRule",
	}
}
