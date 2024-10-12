// Interface for creation of test modes, such as WebhookTestMode.
package ITestMode

type ITestMode interface {
	Run(module string)
}
