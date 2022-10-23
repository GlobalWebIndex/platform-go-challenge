package sql

var (
	chartsUp = `INSERT INTO challenge.charts
	(id,	title, 		x_axis, 	y_axis, 	data) VALUES
	(1,		"chart_1",	"x axis",	"y axis", 	'{"data":1}'),
	(2,		"chart_2",	"x axis",	"y axis", 	'{"data":2}'),
	(3,		"chart_3",	"x axis",	"y axis", 	'{"data":3}');`
	chartsDown = "DELETE FROM  challenge.charts;"

	audiencesUp = `INSERT INTO challenge.audiences
	(id, gender,   country_of_birth, 	age_group, 			hours_spent_online, number_of_purchases_last_month) VALUES
	(1,	 'male',   'gr', 				'young-adults', 	10, 				5),
	(2,	 'female', 'gr', 				'young-adults', 	10, 				5),
	(3,	 'male',   'de', 				'teenagers', 		25, 				2);
	`
	audiencesDown = "DELETE FROM  challenge.audiences;"

	insightsUp = `INSERT INTO challenge.insights
	(id, title) VALUES
	(1,  '40% of millenials spend more than 3hours on social media daily'),
	(2,  '60% of teenagers spend more than 6hours on social media daily'),
	(3,  '10% of seniors spend less than 3hours on social media weekly');
	`
	insightsDown = `DELETE FROM  challenge.insights;`

	dashboardsUp = `INSERT INTO challenge.dashboards
	(id,user_id) VALUES
	(1,	1),
	(2,	2),
	(3,	3);`
	dashboardsDown = "DELETE FROM  challenge.dashboards;"

	d2aUp = `INSERT INTO challenge.dashboards2assets
	(dashboard_id, 	asset_id, 	asset_type, description) VALUES
	(1, 			1, 			'chart', 		'chart description 1'),
	(1, 			2, 			'chart', 		'chart description 2'),
	(1, 			1, 			'audience', 	'audience description 1'),
	(1, 			1, 			'insight', 		'insight description 1'),
	(2, 			1, 			'chart', 		'chart description 1 alt');
	`
	d2aDown = "DELETE FROM  challenge.dashboards2assets;"
)
