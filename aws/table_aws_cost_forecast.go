package aws

import (
	"context"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/costexplorer"

	"github.com/turbot/steampipe-plugin-sdk/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/plugin"
	"github.com/turbot/steampipe-plugin-sdk/plugin/transform"
)

func tableAwsCostForecast(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_cost_forecast",
		Description: "AWS Cost Explorer - Cost Forecast",
		List: &plugin.ListConfig{
			KeyColumns: plugin.SingleColumn("granularity"),
			Hydrate:    listCostForecast,
		},
		Columns: awsColumns([]*plugin.Column{
			{
				Name:        "granularity",
				Description: "",
				Type:        proto.ColumnType_STRING,
				Hydrate:     hydrateCostAndUsageQuals,
			},
			{
				Name:        "period_start",
				Description: "Start timestamp for this cost metric",
				Type:        proto.ColumnType_TIMESTAMP,
				Transform:   transform.FromField("TimePeriod.Start"),
			},
			{
				Name:        "period_end",
				Description: "End timestamp for this cost metric",
				Type:        proto.ColumnType_TIMESTAMP,
				Transform:   transform.FromField("TimePeriod.End"),
			},
			{
				Name:        "mean_value",
				Description: "Average forecasted value",
				Type:        proto.ColumnType_DOUBLE,
			},
		},
		),
	}
}

//// LIST FUNCTION

func listCostForecast(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {

	logger := plugin.Logger(ctx)
	logger.Trace("listCostForecast")

	// Create session
	svc, err := CostExplorerService(ctx, d.ConnectionManager)
	if err != nil {
		return nil, err
	}

	params := buildCostForecastInput(d.KeyColumnQuals)

	output, err := svc.GetCostForecast(params)
	if err != nil {
		logger.Error("listCostForecast", "err", err)
		return nil, err
	}

	// stream the results...
	for _, r := range output.ForecastResultsByTime {
		d.StreamListItem(ctx, r)
	}

	return nil, nil
}

func buildCostForecastInput(keyQuals map[string]*proto.QualValue) *costexplorer.GetCostForecastInput {
	granularity := strings.ToUpper(keyQuals["granularity"].GetStringValue())

	// TO DO - specify metric as qual?   get all cost metrics in parallel?
	//metric := strings.ToUpper(keyQuals["metric"].GetStringValue())
	metric := "UNBLENDED_COST"

	timeFormat := "2006-01-02"
	startTime := time.Now().Format(timeFormat)
	endTime := getForecastEndDateForGranularity(granularity).Format(timeFormat)

	params := &costexplorer.GetCostForecastInput{
		TimePeriod: &costexplorer.DateInterval{
			Start: aws.String(startTime),
			End:   aws.String(endTime),
		},
		Granularity: aws.String(granularity),
		Metric:      aws.String(metric),
	}

	return params
}

func getForecastEndDateForGranularity(granularity string) time.Time {
	switch granularity {
	case "MONTHLY":
		return lastDayOfMonth(12) // 1 year
	case "DAILY":
		return lastDayOfMonth(3) // 3 months
	}
	return lastDayOfMonth(12) // 1 year
}

func lastDayOfMonth(numMonths int) time.Time {
	today := time.Now()
	goneDaysOfMonth := today.Day()

	if goneDaysOfMonth == 1 {
		return today.AddDate(0, numMonths, 0)
	}
	return today.AddDate(0, numMonths+1, -goneDaysOfMonth+1)
}