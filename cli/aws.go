package cli

import (
	"fmt"
	"log"
	"os"

	"github.com/BishopFox/cloudfox/aws"
	"github.com/BishopFox/cloudfox/utils"
	"github.com/aws/aws-sdk-go-v2/service/apigateway"
	"github.com/aws/aws-sdk-go-v2/service/apigatewayv2"
	"github.com/aws/aws-sdk-go-v2/service/apprunner"
	"github.com/aws/aws-sdk-go-v2/service/cloudformation"
	"github.com/aws/aws-sdk-go-v2/service/cloudfront"
	"github.com/aws/aws-sdk-go-v2/service/cloudtrail"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ecr"
	"github.com/aws/aws-sdk-go-v2/service/ecs"
	"github.com/aws/aws-sdk-go-v2/service/efs"
	"github.com/aws/aws-sdk-go-v2/service/eks"
	"github.com/aws/aws-sdk-go-v2/service/elasticloadbalancing"
	"github.com/aws/aws-sdk-go-v2/service/elasticloadbalancingv2"
	"github.com/aws/aws-sdk-go-v2/service/fsx"
	"github.com/aws/aws-sdk-go-v2/service/grafana"
	"github.com/aws/aws-sdk-go-v2/service/iam"
	"github.com/aws/aws-sdk-go-v2/service/lambda"
	"github.com/aws/aws-sdk-go-v2/service/lightsail"
	"github.com/aws/aws-sdk-go-v2/service/mq"
	"github.com/aws/aws-sdk-go-v2/service/opensearch"
	"github.com/aws/aws-sdk-go-v2/service/rds"
	"github.com/aws/aws-sdk-go-v2/service/redshift"
	"github.com/aws/aws-sdk-go-v2/service/route53"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/sagemaker"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
	"github.com/aws/aws-sdk-go-v2/service/ssm"
	"github.com/fatih/color"
	"github.com/kyokomi/emoji"
	"github.com/spf13/cobra"
)

var (
	AWSRegions = []string{"us-east-1", "us-east-2", "us-west-1", "us-west-2", "af-south-1", "ap-east-1", "ap-south-1", "ap-northeast-3", "ap-northeast-2", "ap-southeast-1", "ap-southeast-2", "ap-northeast-1", "ca-central-1", "eu-central-1", "eu-west-1", "eu-west-2", "eu-south-1", "eu-west-3", "eu-north-1", "me-south-1", "sa-east-1"}
	cyan       = color.New(color.FgCyan).SprintFunc()
	green      = color.New(color.FgGreen).SprintFunc()

	AWSProfile         string
	AWSProfilesList    string
	AWSAllProfiles     bool
	AWSProfiles        []string
	AWSOutputFormat    string
	AWSOutputDirectory string
	Verbosity          int
	AWSCommands        = &cobra.Command{
		Use:   "aws",
		Short: "See \"Available Commands\" for AWS Modules",
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Help()
		},
	}

	// The filter is set to "all" when the flag "--filter" is not used
	RoleTrustFilter  string
	RoleTrustCommand = &cobra.Command{
		Use:     "role-trusts",
		Aliases: []string{"roletrusts", "role-trust"},
		Short:   "Enumerate all role trusts",
		Long: "\nUse case examples:\n" +
			"Map all role trusts for caller's account:\n" +
			os.Args[0] + " aws role-trusts\n",
		PreRun: func(cmd *cobra.Command, args []string) {
			for _, profile := range AWSProfiles {
				var caller = utils.AWSWhoami(profile)
				fmt.Printf("[%s] AWS Caller Identity: %s\n", cyan(emoji.Sprintf(":fox:cloudfox v%s :fox:", cmd.Root().Version)), *caller.Arn)
			}
		},
		Run: func(cmd *cobra.Command, args []string) {
			for _, profile := range AWSProfiles {
				m := aws.RoleTrustsModule{
					IAMClient:  iam.NewFromConfig(utils.AWSConfigFileLoader(profile)),
					Caller:     utils.AWSWhoami(profile),
					AWSProfile: profile}
				m.PrintRoleTrusts(AWSOutputFormat, AWSOutputDirectory, Verbosity)
			}
		},
	}

	AccessKeysFilter  string
	AccessKeysCommand = &cobra.Command{
		Use:     "access-keys",
		Aliases: []string{"accesskeys", "keys"},
		Short:   "Enumerate active access keys for all users",
		Long: "\nUse case examples:\n" +
			"Map active access keys:\n" +
			os.Args[0] + " aws access-keys --profile test_account" +
			os.Args[0] + " aws access-keys --filter access_key_id --profile readonly_profile",
		PreRun: func(cmd *cobra.Command, args []string) {
			for _, profile := range AWSProfiles {
				var caller = utils.AWSWhoami(profile)
				fmt.Printf("[%s] AWS Caller Identity: %s\n", cyan(emoji.Sprintf(":fox:cloudfox v%s :fox:", cmd.Root().Version)), *caller.Arn)
			}
		},
		Run: func(cmd *cobra.Command, args []string) {
			for _, profile := range AWSProfiles {
				m := aws.AccessKeysModule{
					IAMClient:  iam.NewFromConfig(utils.AWSConfigFileLoader(profile)),
					Caller:     utils.AWSWhoami(profile),
					AWSProfile: profile,
				}
				m.PrintAccessKeys(AccessKeysFilter, AWSOutputFormat, AWSOutputDirectory, Verbosity)
			}
		},
	}

	BucketsCommand = &cobra.Command{
		Use:     "buckets",
		Aliases: []string{"bucket"},
		Short:   "Enumerate all of the buckets. Get loot file with s3 commands to list/download bucket contents",
		Long: "\nUse case examples:\n" +
			"List all buckets create a file with pre-populated aws s3 commands:\n" +
			os.Args[0] + " aws buckets --profile test_account",
		PreRun: func(cmd *cobra.Command, args []string) {
			for _, profile := range AWSProfiles {
				var caller = utils.AWSWhoami(profile)
				fmt.Printf("[%s] AWS Caller Identity: %s\n", cyan(emoji.Sprintf(":fox:cloudfox v%s :fox:", cmd.Root().Version)), *caller.Arn)
			}
		},
		Run: func(cmd *cobra.Command, args []string) {
			for _, profile := range AWSProfiles {
				m := aws.BucketsModule{
					S3Client:   s3.NewFromConfig(utils.AWSConfigFileLoader(profile)),
					Caller:     utils.AWSWhoami(profile),
					AWSProfile: profile,
				}
				m.PrintBuckets(AWSOutputFormat, AWSOutputDirectory, Verbosity)
			}
		},
	}

	// This filter could be an instance ID or a TXT file with instance IDs separated by a new line.
	InstancesFilter                   string
	InstanceMapUserDataAttributesOnly bool
	InstancesCommand                  = &cobra.Command{
		Use:     "instances",
		Aliases: []string{"instance"},
		Short:   "Enumerate all instances along with assigned IPs, profiles, and user-data",
		Long: "\nUse case examples:\n" +
			os.Args[0] + " aws instances --profile readonly_profile",
		PreRun: func(cmd *cobra.Command, args []string) {
			for _, profile := range AWSProfiles {
				var caller = utils.AWSWhoami(profile)
				fmt.Printf("[%s] AWS Caller Identity: %s\n", cyan(emoji.Sprintf(":fox:cloudfox v%s :fox:", cmd.Root().Version)), *caller.Arn)
			}
		},
		Run: func(cmd *cobra.Command, args []string) {
			for _, profile := range AWSProfiles {
				m := aws.InstancesModule{
					EC2Client:  ec2.NewFromConfig(utils.AWSConfigFileLoader(profile)),
					Caller:     utils.AWSWhoami(profile),
					AWSRegions: AWSRegions,

					UserDataAttributesOnly: InstanceMapUserDataAttributesOnly,
					AWSProfile:             profile,
				}
				m.Instances(InstancesFilter, AWSOutputFormat, AWSOutputDirectory, Verbosity)
			}
		},
	}

	InventoryCommand = &cobra.Command{
		Use:   "inventory",
		Short: "Gain a rough understanding of size of the account and preferred regions",
		Long: "\nUse case examples:\n" +
			os.Args[0] + " aws inventory --profile readonly_profile",
		PreRun: func(cmd *cobra.Command, args []string) {
			for _, profile := range AWSProfiles {
				var caller = utils.AWSWhoami(profile)
				fmt.Printf("[%s] AWS Caller Identity: %s\n", cyan(emoji.Sprintf(":fox:cloudfox v%s :fox:", cmd.Root().Version)), *caller.Arn)
			}
		},
		Run: func(cmd *cobra.Command, args []string) {
			for _, profile := range AWSProfiles {
				m := aws.Inventory2Module{
					EC2Client:            ec2.NewFromConfig(utils.AWSConfigFileLoader(profile)),
					ECSClient:            ecs.NewFromConfig(utils.AWSConfigFileLoader(profile)),
					EKSClient:            eks.NewFromConfig(utils.AWSConfigFileLoader(profile)),
					S3Client:             s3.NewFromConfig(utils.AWSConfigFileLoader(profile)),
					LambdaClient:         lambda.NewFromConfig(utils.AWSConfigFileLoader(profile)),
					CloudFormationClient: cloudformation.NewFromConfig(utils.AWSConfigFileLoader(profile)),
					SecretsManagerClient: secretsmanager.NewFromConfig(utils.AWSConfigFileLoader(profile)),
					SSMClient:            ssm.NewFromConfig(utils.AWSConfigFileLoader(profile)),
					RDSClient:            rds.NewFromConfig(utils.AWSConfigFileLoader(profile)),
					APIGatewayv2Client:   apigatewayv2.NewFromConfig(utils.AWSConfigFileLoader(profile)),
					ELBClient:            elasticloadbalancing.NewFromConfig(utils.AWSConfigFileLoader(profile)),
					ELBv2Client:          elasticloadbalancingv2.NewFromConfig(utils.AWSConfigFileLoader(profile)),
					IAMClient:            iam.NewFromConfig(utils.AWSConfigFileLoader(profile)),
					MQClient:             mq.NewFromConfig(utils.AWSConfigFileLoader(profile)),
					OpenSearchClient:     opensearch.NewFromConfig(utils.AWSConfigFileLoader(profile)),
					GrafanaClient:        grafana.NewFromConfig(utils.AWSConfigFileLoader(profile)),
					APIGatewayClient:     apigateway.NewFromConfig(utils.AWSConfigFileLoader(profile)),
					RedshiftClient:       redshift.NewFromConfig(utils.AWSConfigFileLoader(profile)),
					CloudfrontClient:     cloudfront.NewFromConfig(utils.AWSConfigFileLoader(profile)),
					AppRunnerClient:      apprunner.NewFromConfig(utils.AWSConfigFileLoader(profile)),
					LightsailClient:      lightsail.NewFromConfig(utils.AWSConfigFileLoader(profile)),

					Caller:     utils.AWSWhoami(profile),
					AWSRegions: AWSRegions,
					AWSProfile: profile,
				}
				m.PrintInventoryPerRegion(AWSOutputFormat, AWSOutputDirectory, Verbosity)
			}
		},
	}

	EndpointsCommand = &cobra.Command{
		Use:     "endpoints",
		Aliases: []string{"endpoint"},
		Short:   "Enumerates endpoints from various services. Get a loot file with http endpoints to scan.",
		Long: "\nUse case examples:\n" +
			os.Args[0] + " aws endpoints --profile readonly_profile",
		PreRun: func(cmd *cobra.Command, args []string) {
			for _, profile := range AWSProfiles {
				var caller = utils.AWSWhoami(profile)
				fmt.Printf("[%s] AWS Caller Identity: %s\n", cyan(emoji.Sprintf(":fox:cloudfox v%s :fox:", cmd.Root().Version)), *caller.Arn)
			}
		},
		Run: func(cmd *cobra.Command, args []string) {
			for _, profile := range AWSProfiles {
				m := aws.EndpointsModule{
					EKSClient:          eks.NewFromConfig(utils.AWSConfigFileLoader(profile)),
					LambdaClient:       lambda.NewFromConfig(utils.AWSConfigFileLoader(profile)),
					MQClient:           mq.NewFromConfig(utils.AWSConfigFileLoader(profile)),
					OpenSearchClient:   opensearch.NewFromConfig(utils.AWSConfigFileLoader(profile)),
					GrafanaClient:      grafana.NewFromConfig(utils.AWSConfigFileLoader(profile)),
					ELBClient:          elasticloadbalancing.NewFromConfig(utils.AWSConfigFileLoader(profile)),
					APIGatewayClient:   apigateway.NewFromConfig(utils.AWSConfigFileLoader(profile)),
					ELBv2Client:        elasticloadbalancingv2.NewFromConfig(utils.AWSConfigFileLoader(profile)),
					APIGatewayv2Client: apigatewayv2.NewFromConfig(utils.AWSConfigFileLoader(profile)),
					RDSClient:          rds.NewFromConfig(utils.AWSConfigFileLoader(profile)),
					RedshiftClient:     redshift.NewFromConfig(utils.AWSConfigFileLoader(profile)),
					S3Client:           s3.NewFromConfig(utils.AWSConfigFileLoader(profile)),
					CloudfrontClient:   cloudfront.NewFromConfig(utils.AWSConfigFileLoader(profile)),
					AppRunnerClient:    apprunner.NewFromConfig(utils.AWSConfigFileLoader(profile)),
					LightsailClient:    lightsail.NewFromConfig(utils.AWSConfigFileLoader(profile)),

					Caller:     utils.AWSWhoami(profile),
					AWSRegions: AWSRegions,
					AWSProfile: profile,
				}
				m.PrintEndpoints(AWSOutputFormat, AWSOutputDirectory, Verbosity)
			}
		},
	}

	SecretsCommand = &cobra.Command{
		Use:     "secrets",
		Aliases: []string{"secret"},
		Short:   "Enumerate secrets from secrets manager and SSM",
		Long: "\nUse case examples:\n" +
			os.Args[0] + " aws secrets --profile readonly_profile",
		PreRun: func(cmd *cobra.Command, args []string) {
			for _, profile := range AWSProfiles {
				var caller = utils.AWSWhoami(profile)
				fmt.Printf("[%s] AWS Caller Identity: %s\n", cyan(emoji.Sprintf(":fox:cloudfox v%s :fox:", cmd.Root().Version)), *caller.Arn)
			}
		},
		Run: func(cmd *cobra.Command, args []string) {
			for _, profile := range AWSProfiles {
				m := aws.SecretsModule{
					SecretsManagerClient: secretsmanager.NewFromConfig(utils.AWSConfigFileLoader(profile)),
					SSMClient:            ssm.NewFromConfig(utils.AWSConfigFileLoader(profile)),

					Caller:     utils.AWSWhoami(profile),
					AWSRegions: AWSRegions,
					AWSProfile: profile,
				}
				m.PrintSecrets(AWSOutputFormat, AWSOutputDirectory, Verbosity)
			}
		},
	}

	Route53Command = &cobra.Command{
		Use:     "route53",
		Aliases: []string{"dns", "route", "routes"},
		Short:   "Enumerate all records from all zones managed by route53. Get a loot file with A records you can scan",
		Long: "\nUse case examples:\n" +
			os.Args[0] + " aws route53 --profile readonly_profile",
		PreRun: func(cmd *cobra.Command, args []string) {
			for _, profile := range AWSProfiles {
				var caller = utils.AWSWhoami(profile)
				fmt.Printf("[%s] AWS Caller Identity: %s\n", cyan(emoji.Sprintf(":fox:cloudfox v%s :fox:", cmd.Root().Version)), *caller.Arn)
			}
		},
		Run: func(cmd *cobra.Command, args []string) {
			for _, profile := range AWSProfiles {
				m := aws.Route53Module{
					Route53Client: route53.NewFromConfig(utils.AWSConfigFileLoader(profile)),

					Caller:     utils.AWSWhoami(profile),
					AWSRegions: AWSRegions,
					AWSProfile: profile,
				}
				m.PrintRoute53(AWSOutputFormat, AWSOutputDirectory, Verbosity)
			}
		},
	}

	ECRCommand = &cobra.Command{
		Use:     "ecr",
		Aliases: []string{"repos", "repo", "repositories"},
		Short:   "Enumerate the most recently pushed image URI from all repositories. Get a loot file with commands to pull images",
		Long: "\nUse case examples:\n" +
			os.Args[0] + " aws ecr --profile readonly_profile",
		PreRun: func(cmd *cobra.Command, args []string) {
			for _, profile := range AWSProfiles {
				var caller = utils.AWSWhoami(profile)
				fmt.Printf("[%s] AWS Caller Identity: %s\n", cyan(emoji.Sprintf(":fox:cloudfox v%s :fox:", cmd.Root().Version)), *caller.Arn)
			}
		},
		Run: func(cmd *cobra.Command, args []string) {
			for _, profile := range AWSProfiles {
				m := aws.ECRModule{
					ECRClient: ecr.NewFromConfig(utils.AWSConfigFileLoader(profile)),

					Caller:     utils.AWSWhoami(profile),
					AWSRegions: AWSRegions,
					AWSProfile: profile,
				}
				m.PrintECR(AWSOutputFormat, AWSOutputDirectory, Verbosity)
			}
		},
	}
	OutboundAssumedRolesDays    int
	OutboundAssumedRolesCommand = &cobra.Command{
		Use:     "outbound-assumed-roles",
		Aliases: []string{"assumedroles", "assumeroles", "outboundassumedroles"},
		Short:   "Find the roles that have been assumed by principals in this account",
		Long: "\nUse case examples:\n" +
			os.Args[0] + " aws outbound-assumed-roles --profile readonly_profile",
		PreRun: func(cmd *cobra.Command, args []string) {
			for _, profile := range AWSProfiles {
				var caller = utils.AWSWhoami(profile)
				fmt.Printf("[%s] AWS Caller Identity: %s\n", cyan(emoji.Sprintf(":fox:cloudfox v%s :fox:", cmd.Root().Version)), *caller.Arn)
			}
		},
		Run: func(cmd *cobra.Command, args []string) {
			for _, profile := range AWSProfiles {
				m := aws.OutboundAssumedRolesModule{
					CloudTrailClient: cloudtrail.NewFromConfig(utils.AWSConfigFileLoader(profile)),

					Caller:     utils.AWSWhoami(profile),
					AWSRegions: AWSRegions,
					AWSProfile: profile,
				}
				m.PrintOutboundRoleTrusts(OutboundAssumedRolesDays, AWSOutputFormat, AWSOutputDirectory, Verbosity)
			}
		},
	}

	EnvsCommand = &cobra.Command{
		Use:     "env-vars",
		Aliases: []string{"envs", "envvars", "env"},
		Short:   "Enumerate the environment variables from mutliple services that have them",
		Long: "\nUse case examples:\n" +
			os.Args[0] + " aws env-vars --profile readonly_profile",
		PreRun: func(cmd *cobra.Command, args []string) {
			for _, profile := range AWSProfiles {
				var caller = utils.AWSWhoami(profile)
				fmt.Printf("[%s] AWS Caller Identity: %s\n", cyan(emoji.Sprintf(":fox:cloudfox v%s :fox:", cmd.Root().Version)), *caller.Arn)
			}
		},
		Run: func(cmd *cobra.Command, args []string) {
			for _, profile := range AWSProfiles {
				m := aws.EnvsModule{

					Caller:          utils.AWSWhoami(profile),
					AWSRegions:      AWSRegions,
					AWSProfile:      profile,
					ECSClient:       ecs.NewFromConfig(utils.AWSConfigFileLoader(profile)),
					AppRunnerClient: apprunner.NewFromConfig(utils.AWSConfigFileLoader(profile)),
					LambdaClient:    lambda.NewFromConfig(utils.AWSConfigFileLoader(profile)),
					LightsailClient: lightsail.NewFromConfig(utils.AWSConfigFileLoader(profile)),
					SagemakerClient: sagemaker.NewFromConfig(utils.AWSConfigFileLoader(profile)),
				}
				m.PrintEnvs(AWSOutputFormat, AWSOutputDirectory, Verbosity)
			}
		},
	}

	PrincipalsCommand = &cobra.Command{
		Use:     "principals",
		Aliases: []string{"principal"},
		Short:   "Enumerate IAM users and Roles so you have the data at your fingertips",
		Long: "\nUse case examples:\n" +
			os.Args[0] + " aws principals --profile readonly_profile",
		PreRun: func(cmd *cobra.Command, args []string) {
			for _, profile := range AWSProfiles {
				var caller = utils.AWSWhoami(profile)
				fmt.Printf("[%s] AWS Caller Identity: %s\n", cyan(emoji.Sprintf(":fox:cloudfox v%s :fox:", cmd.Root().Version)), *caller.Arn)
			}
		},
		Run: func(cmd *cobra.Command, args []string) {
			for _, profile := range AWSProfiles {
				m := aws.IamPrincipalsModule{
					IAMClient:  iam.NewFromConfig(utils.AWSConfigFileLoader(profile)),
					Caller:     utils.AWSWhoami(profile),
					AWSProfile: profile,
				}
				m.PrintIamPrincipals(AWSOutputFormat, AWSOutputDirectory, Verbosity)
			}
		},
	}

	PermissionsPrincipal string
	PermissionsCommand   = &cobra.Command{
		Use:     "permissions",
		Aliases: []string{"perms", "permission"},
		Short:   "Enumerate IAM permissions per principal",
		Long: "\nUse case examples:\n" +
			os.Args[0] + " aws permissions --profile profile\n" +
			os.Args[0] + " aws permissions --profile profile --principal arn:aws:iam::111111111111:role/test123",
		PreRun: func(cmd *cobra.Command, args []string) {
			for _, profile := range AWSProfiles {
				var caller = utils.AWSWhoami(profile)
				fmt.Printf("[%s] AWS Caller Identity: %s\n", cyan(emoji.Sprintf(":fox:cloudfox v%s :fox:", cmd.Root().Version)), *caller.Arn)
			}
		},
		Run: func(cmd *cobra.Command, args []string) {
			for _, profile := range AWSProfiles {
				m := aws.IamPermissionsModule{
					IAMClient:  iam.NewFromConfig(utils.AWSConfigFileLoader(profile)),
					Caller:     utils.AWSWhoami(profile),
					AWSProfile: profile,
				}
				m.PrintIamPermissions(AWSOutputFormat, AWSOutputDirectory, Verbosity, PermissionsPrincipal)
			}
		},
	}

	SimulatorResource   string
	SimulatorAction     string
	SimulatorPrincipal  string
	IamSimulatorCommand = &cobra.Command{
		Use:     "iam-simulator",
		Aliases: []string{"iamsimulator", "simulator"},
		Short:   "Wrapper around the AWS IAM Simulate Principal Policy command",
		Long: "\nUse case examples:\n" +
			os.Args[0] + " aws iam-simulator --profile readonly_profile",
		PreRun: func(cmd *cobra.Command, args []string) {
			for _, profile := range AWSProfiles {
				var caller = utils.AWSWhoami(profile)
				fmt.Printf("[%s] AWS Caller Identity: %s\n", cyan(emoji.Sprintf(":fox:cloudfox v%s :fox:", cmd.Root().Version)), *caller.Arn)
			}
		},
		Run: func(cmd *cobra.Command, args []string) {
			for _, profile := range AWSProfiles {
				m := aws.IamSimulatorModule{
					IAMClient: iam.NewFromConfig(utils.AWSConfigFileLoader(profile)),

					Caller:     utils.AWSWhoami(profile),
					AWSProfile: profile,
				}
				m.PrintIamSimulator(SimulatorPrincipal, SimulatorAction, SimulatorResource, AWSOutputFormat, AWSOutputDirectory, Verbosity)
			}
		},
	}

	FilesystemsCommand = &cobra.Command{
		Use:     "filesystems",
		Aliases: []string{"filesystem"},
		Short:   "Enumerate the EFS and FSx filesystems. Get a loot file with mount commands",
		Long: "\nUse case examples:\n" +
			os.Args[0] + " aws filesystems --profile readonly_profile",
		PreRun: func(cmd *cobra.Command, args []string) {
			for _, profile := range AWSProfiles {
				var caller = utils.AWSWhoami(profile)
				fmt.Printf("[%s] AWS Caller Identity: %s\n", cyan(emoji.Sprintf(":fox:cloudfox v%s :fox:", cmd.Root().Version)), *caller.Arn)
			}
		},
		Run: func(cmd *cobra.Command, args []string) {
			for _, profile := range AWSProfiles {
				filesystems := aws.FilesystemsModule{
					EFSClient: efs.NewFromConfig(utils.AWSConfigFileLoader(profile)),
					FSxClient: fsx.NewFromConfig(utils.AWSConfigFileLoader(profile)),

					Caller:     utils.AWSWhoami(profile),
					AWSProfile: profile,
					AWSRegions: AWSRegions,
				}
				filesystems.PrintFilesystems(AWSOutputFormat, AWSOutputDirectory, Verbosity)
			}
		},
	}

	// RAMCommand = &cobra.Command{
	// 	Use:   "ram",
	// 	Short: "Enumerate cross-account shared resources",
	// 	Long: "\nUse case examples:\n" +
	// 		os.Args[0] + " aws ram --profile readonly_profile",
	// 	PreRun: func(cmd *cobra.Command, args []string) {
	// 		var caller = utils.AWSWhoami(AWSProfile)
	// 		fmt.Printf("[%s] AWS Caller Identity: %s\n", cyan(emoji.Sprintf(":fox:cloudfox v%s :fox:", cmd.Root().Version)), *caller.Arn)
	// 	},
	// 	Run: func(cmd *cobra.Command, args []string) {
	// 		m := aws.RAMModule{
	// 			RAMClient: ram.NewFromConfig(utils.AWSConfigFileLoader(AWSProfile)),

	// 			Caller:     utils.AWSWhoami(AWSProfile),
	// 			AWSRegions: AWSRegions,
	// 			AWSProfile: AWSProfile,
	// 		}
	// 		m.PrintRAM(AWSOutputFormat, AWSOutputDirectory, Verbosity)
	// 	},
	// }

	AllChecksCommand = &cobra.Command{

		Use:     "all-checks",
		Aliases: []string{"allchecks", "all"},
		Short:   "Run all of the other checks (excluding outbound-assumed-roles)",
		Long: "\nUse case examples:\n" +
			os.Args[0] + " aws all-checks --profile readonly_profile", //TODO add examples? os.Args[0] + " aws all-checks --profiles profiles.txt, os.Args[0] + " aws all-checks --all-profiles""
		PreRun: func(cmd *cobra.Command, args []string) {
			for _, profile := range AWSProfiles {
				var caller = utils.AWSWhoami(profile)
				fmt.Printf("[%s] AWS Caller Identity: %s\n", cyan(emoji.Sprintf(":fox:cloudfox v%s :fox:", cmd.Root().Version)), *caller.Arn)
			}
		},
		Run: func(cmd *cobra.Command, args []string) {
			for _, profile := range AWSProfiles {
				ec2Client := ec2.NewFromConfig(utils.AWSConfigFileLoader(profile))
				eksClient := eks.NewFromConfig(utils.AWSConfigFileLoader(profile))
				s3Client := s3.NewFromConfig(utils.AWSConfigFileLoader(profile))
				lambdaClient := lambda.NewFromConfig(utils.AWSConfigFileLoader(profile))
				cloudFormationClient := cloudformation.NewFromConfig(utils.AWSConfigFileLoader(profile))
				secretsManagerClient := secretsmanager.NewFromConfig(utils.AWSConfigFileLoader(profile))
				rdsClient := rds.NewFromConfig(utils.AWSConfigFileLoader(profile))
				apiGatewayv2Client := apigatewayv2.NewFromConfig(utils.AWSConfigFileLoader(profile))
				apiGatewayClient := apigateway.NewFromConfig(utils.AWSConfigFileLoader(profile))
				elbClient := elasticloadbalancing.NewFromConfig(utils.AWSConfigFileLoader(profile))
				elbv2Client := elasticloadbalancingv2.NewFromConfig(utils.AWSConfigFileLoader(profile))
				iamClient := iam.NewFromConfig(utils.AWSConfigFileLoader(profile))
				mqClient := mq.NewFromConfig(utils.AWSConfigFileLoader(profile))
				openSearchClient := opensearch.NewFromConfig(utils.AWSConfigFileLoader(profile))
				grafanaClient := grafana.NewFromConfig(utils.AWSConfigFileLoader(profile))
				redshiftClient := redshift.NewFromConfig(utils.AWSConfigFileLoader(profile))
				cloudfrontClient := cloudfront.NewFromConfig(utils.AWSConfigFileLoader(profile))
				appRunnerClient := apprunner.NewFromConfig(utils.AWSConfigFileLoader(profile))
				lightsailClient := lightsail.NewFromConfig(utils.AWSConfigFileLoader(profile))
				route53Client := route53.NewFromConfig(utils.AWSConfigFileLoader(profile))
				efsClient := efs.NewFromConfig(utils.AWSConfigFileLoader(profile))
				fsxClient := fsx.NewFromConfig(utils.AWSConfigFileLoader(profile))
				ecsClient := ecs.NewFromConfig(utils.AWSConfigFileLoader(profile))
				sagemakerClient := sagemaker.NewFromConfig(utils.AWSConfigFileLoader(profile))
				ecrClient := ecr.NewFromConfig(utils.AWSConfigFileLoader(profile))
				ssmClient := ssm.NewFromConfig(utils.AWSConfigFileLoader(profile))

				fmt.Printf("[%s] %s\n", cyan(emoji.Sprintf(":fox:cloudfox :fox:")), green("Getting a lay of the land, aka \"What regions is this account using?\""))
				inventory2 := aws.Inventory2Module{
					EC2Client:            ec2Client,
					ECSClient:            ecsClient,
					EKSClient:            eksClient,
					S3Client:             s3Client,
					LambdaClient:         lambdaClient,
					CloudFormationClient: cloudFormationClient,
					SecretsManagerClient: secretsManagerClient,
					SSMClient:            ssmClient,
					RDSClient:            rdsClient,
					APIGatewayv2Client:   apiGatewayv2Client,
					APIGatewayClient:     apiGatewayClient,
					ELBClient:            elbClient,
					ELBv2Client:          elbv2Client,
					IAMClient:            iamClient,
					MQClient:             mqClient,
					OpenSearchClient:     openSearchClient,
					GrafanaClient:        grafanaClient,
					RedshiftClient:       redshiftClient,
					CloudfrontClient:     cloudfrontClient,
					AppRunnerClient:      appRunnerClient,
					LightsailClient:      lightsailClient,

					Caller:     utils.AWSWhoami(profile),
					AWSRegions: AWSRegions,
					AWSProfile: profile,
				}
				inventory2.PrintInventoryPerRegion(AWSOutputFormat, AWSOutputDirectory, Verbosity)
				//time.Sleep(time.Second * 5)
				// Service and endpoint enum section
				fmt.Printf("[%s] %s\n", cyan(emoji.Sprintf(":fox:cloudfox :fox:")), green("Gathering the info you'll want for your application & service enumeration needs."))

				instances := aws.InstancesModule{
					EC2Client:  ec2Client,
					Caller:     utils.AWSWhoami(profile),
					AWSRegions: AWSRegions,

					UserDataAttributesOnly: false,
					AWSProfile:             profile,
				}
				instances.Instances(InstancesFilter, AWSOutputFormat, AWSOutputDirectory, Verbosity)
				route53 := aws.Route53Module{
					Route53Client: route53Client,

					Caller:     utils.AWSWhoami(profile),
					AWSRegions: AWSRegions,
					AWSProfile: profile,
				}
				route53.PrintRoute53(AWSOutputFormat, AWSOutputDirectory, Verbosity)

				filesystems := aws.FilesystemsModule{
					EFSClient:  efsClient,
					FSxClient:  fsxClient,
					Caller:     utils.AWSWhoami(profile),
					AWSProfile: profile,
					AWSRegions: AWSRegions,
				}
				filesystems.PrintFilesystems(AWSOutputFormat, AWSOutputDirectory, Verbosity)

				endpoints := aws.EndpointsModule{

					EKSClient:          eksClient,
					S3Client:           s3Client,
					LambdaClient:       lambdaClient,
					RDSClient:          rdsClient,
					APIGatewayv2Client: apiGatewayv2Client,
					APIGatewayClient:   apiGatewayClient,
					ELBClient:          elbClient,
					ELBv2Client:        elbv2Client,
					MQClient:           mqClient,
					OpenSearchClient:   openSearchClient,
					GrafanaClient:      grafanaClient,
					RedshiftClient:     redshiftClient,
					CloudfrontClient:   cloudfrontClient,
					AppRunnerClient:    appRunnerClient,
					LightsailClient:    lightsailClient,

					Caller:     utils.AWSWhoami(profile),
					AWSRegions: AWSRegions,
					AWSProfile: profile,
				}

				endpoints.PrintEndpoints(AWSOutputFormat, AWSOutputDirectory, Verbosity)
				// Secrets section
				fmt.Printf("[%s] %s\n", cyan(emoji.Sprintf(":fox:cloudfox :fox:")), green("Looking for secrets hidden between the seat cushions."))

				ec2UserData := aws.InstancesModule{
					EC2Client:  ec2Client,
					Caller:     utils.AWSWhoami(profile),
					AWSRegions: AWSRegions,

					UserDataAttributesOnly: true,
					AWSProfile:             profile,
				}
				ec2UserData.Instances(InstancesFilter, AWSOutputFormat, AWSOutputDirectory, Verbosity)
				envsMod := aws.EnvsModule{

					Caller:          utils.AWSWhoami(profile),
					AWSRegions:      AWSRegions,
					AWSProfile:      profile,
					ECSClient:       ecsClient,
					AppRunnerClient: appRunnerClient,
					LambdaClient:    lambdaClient,
					LightsailClient: lightsailClient,
					SagemakerClient: sagemakerClient,
				}
				envsMod.PrintEnvs(AWSOutputFormat, AWSOutputDirectory, Verbosity)

				// CPT Enum
				//fmt.Printf("[%s] %s\n", cyan(emoji.Sprintf(":fox:cloudfox :fox:")), green("Gathering some other info that is often useful."))
				fmt.Printf("[%s] %s\n", cyan(emoji.Sprintf(":fox:cloudfox :fox:")), green("Arming you with the data you'll need for privesc quests."))

				buckets := aws.BucketsModule{
					S3Client:   s3Client,
					Caller:     utils.AWSWhoami(profile),
					AWSProfile: profile,
				}
				buckets.PrintBuckets(AWSOutputFormat, AWSOutputDirectory, Verbosity)

				ecr := aws.ECRModule{
					ECRClient:  ecrClient,
					Caller:     utils.AWSWhoami(profile),
					AWSRegions: AWSRegions,
					AWSProfile: profile,
				}
				ecr.PrintECR(AWSOutputFormat, AWSOutputDirectory, Verbosity)

				secrets := aws.SecretsModule{
					SecretsManagerClient: secretsManagerClient,
					SSMClient:            ssmClient,

					Caller:     utils.AWSWhoami(profile),
					AWSRegions: AWSRegions,
					AWSProfile: profile,
				}
				secrets.PrintSecrets(AWSOutputFormat, AWSOutputDirectory, Verbosity)

				// IAM privesc section
				fmt.Printf("[%s] %s\n", cyan(emoji.Sprintf(":fox:cloudfox :fox:")), green("IAM is complicated. Complicated usually means misconfigurations. You'll want to pay attention here."))
				principals := aws.IamPrincipalsModule{
					IAMClient:  iamClient,
					Caller:     utils.AWSWhoami(profile),
					AWSProfile: profile,
				}
				principals.PrintIamPrincipals(AWSOutputFormat, AWSOutputDirectory, Verbosity)
				permissions := aws.IamPermissionsModule{
					IAMClient:  iamClient,
					Caller:     utils.AWSWhoami(profile),
					AWSProfile: profile,
				}
				permissions.PrintIamPermissions(AWSOutputFormat, AWSOutputDirectory, Verbosity, PermissionsPrincipal)
				accessKeys := aws.AccessKeysModule{
					IAMClient:  iam.NewFromConfig(utils.AWSConfigFileLoader(profile)),
					Caller:     utils.AWSWhoami(profile),
					AWSProfile: profile,
				}
				accessKeys.PrintAccessKeys(AccessKeysFilter, AWSOutputFormat, AWSOutputDirectory, Verbosity)
				inboundRoleTrusts := aws.RoleTrustsModule{
					IAMClient:  iamClient,
					Caller:     utils.AWSWhoami(profile),
					AWSProfile: profile,
				}
				inboundRoleTrusts.PrintRoleTrusts(AWSOutputFormat, AWSOutputDirectory, Verbosity)
				iamSimulator := aws.IamSimulatorModule{
					IAMClient:  iamClient,
					Caller:     utils.AWSWhoami(profile),
					AWSProfile: profile,
				}
				iamSimulator.PrintIamSimulator(SimulatorPrincipal, SimulatorAction, SimulatorResource, AWSOutputFormat, AWSOutputDirectory, Verbosity)

				fmt.Printf("[%s] %s\n", cyan(emoji.Sprintf(":fox:cloudfox :fox:")), green("That's it! Check your output files for situational awareness and check your loot files for next steps."))
				fmt.Printf("[%s] %s\n\n", cyan(emoji.Sprintf(":fox:cloudfox :fox:")), green("FYI, we skipped the outbound-assumed-roles module in all-checks (really long run time). Make sure to try it out manually."))
			}
		},
	}
)

func init() {
	cobra.OnInitialize(initAWSProfiles)
	// Principal Trusts Module Flags
	RoleTrustCommand.Flags().StringVarP(&RoleTrustFilter, "filter", "t", "all", "[AccountNumber | PrincipalARN | PrincipalName | ServiceName]")

	// Map Access Keys Module Flags
	AccessKeysCommand.Flags().StringVarP(&AccessKeysFilter, "filter", "t", "none", "Access key ID to search for")

	// IAM Simulator Module Flags
	//IamSimulatorCommand.Flags().StringVarP(&IamSimulatorFilter, "filter", "f", "none", "Access key ID to search for")

	// Instances Map Module Flags
	InstancesCommand.Flags().StringVarP(&InstancesFilter, "filter", "t", "all", "[InstanceID | InstanceIDsFile]")
	InstancesCommand.Flags().BoolVarP(&InstanceMapUserDataAttributesOnly, "userdata", "u", false, "Use this flag to retrieve only the userData attribute from EC2 instances.")

	//  outbound-assumed-roles module flags
	OutboundAssumedRolesCommand.Flags().IntVarP(&OutboundAssumedRolesDays, "days", "d", 7, "How many days of CloudTrail events should we go back and look at.")

	//  iam-simulator module flags
	IamSimulatorCommand.Flags().StringVar(&SimulatorPrincipal, "principal", "", "Principal Arn")
	IamSimulatorCommand.Flags().StringVar(&SimulatorAction, "action", "", "Action")
	IamSimulatorCommand.Flags().StringVar(&SimulatorResource, "resource", "*", "Resource")

	//  iam-simulator module flags
	PermissionsCommand.Flags().StringVar(&PermissionsPrincipal, "principal", "", "Principal Arn")

	// Global flags for the AWS modules
	AWSCommands.PersistentFlags().StringVarP(&AWSProfile, "profile", "p", "", "AWS CLI Profile Name")
	AWSCommands.PersistentFlags().StringVarP(&AWSProfilesList, "profiles-list", "l", "", "List of AWS CLI profile names")
	AWSCommands.PersistentFlags().BoolVarP(&AWSAllProfiles, "all-profiles", "a", false, "Use all AWS CLI profiles in AWS credentials file")
	AWSCommands.PersistentFlags().StringVarP(&AWSOutputFormat, "output", "o", "all", "[\"table\" | \"csv\" | \"all\" ]")
	AWSCommands.PersistentFlags().IntVarP(&Verbosity, "verbosity", "v", 1, "1 = Print control messages only\n2 = Print control messages, module output\n3 = Print control messages, module output, and loot file output\n")
	AWSCommands.PersistentFlags().StringVar(&AWSOutputDirectory, "outdir", ".", "Output Directory ")

	AWSCommands.AddCommand(
		AllChecksCommand,
		RoleTrustCommand,
		AccessKeysCommand,
		InstancesCommand,
		InventoryCommand,
		EndpointsCommand,
		SecretsCommand,
		Route53Command,
		ECRCommand,
		OutboundAssumedRolesCommand,
		EnvsCommand,
		PrincipalsCommand,
		IamSimulatorCommand,
		FilesystemsCommand,
		BucketsCommand,
		PermissionsCommand,
	)

}

func initAWSProfiles() {
	// Ensure only one profile setting is chosen. If
	if AWSProfile != "" && AWSProfilesList != "" || AWSProfile != "" && AWSAllProfiles || AWSProfilesList != "" && AWSAllProfiles {
		log.Fatalf("[-] Error specifying AWS profiles. Choose only one of -p/--profile, -a/--all-profiles, -l/--profiles-list. Exiting...")
	} else if AWSProfile != "" {
		AWSProfiles = append(AWSProfiles, AWSProfile)
	} else if AWSProfilesList != "" {
		// Written like so to enable testing while still being readable
		AWSProfiles = utils.GetSelectedAWSProfiles(AWSProfilesList)
	} else if AWSAllProfiles {
		AWSProfiles = utils.GetAllAWSProfiles()
	} else {
		AWSProfiles = append(AWSProfiles, "default")
	}
}
