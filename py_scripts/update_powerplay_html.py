import os
import json
from bs4 import BeautifulSoup
from selenium import webdriver

from jinja2 import Environment, FileSystemLoader


options = webdriver.ChromeOptions()
options.add_argument("--headless")  # Run in headless mode (no GUI)
driver = webdriver.Chrome(options=options)

# Fetch the HTML content
url = "https://spela.svenskaspel.se/powerplay/"
driver.get(url)
rendered_html = driver.page_source

# Parse the HTML
soup = BeautifulSoup(rendered_html, "html.parser")

# Load template powerplay script
env = Environment(loader=FileSystemLoader(os.path.join(os.getcwd(), "html", "templates")))
p_template = env.get_template('powerplay.html')

teams_home = []
teams_away = []

opt1_percs = []
optX_percs = []
opt2_percs = []

opt1_odds = []
optX_odds = []
opt2_odds = []

# Loop throup each match ('coupon-row') and 
row_containers = soup.find(class_="coupon-rows")
for container in row_containers:
    home_team_html = container.find(class_="participant home-participant")
    away_team_html = container.find(class_="participant away-participant")
    teams_home.append(home_team_html.text)
    teams_away.append(away_team_html.text)

    stats_info_html = container.find_all(class_="stat-trend stat-trend-neutral")
    
    percs_html = stats_info_html[0]
    perc1 = stats_info_html[0]
    percX = stats_info_html[1]
    perc2 = stats_info_html[2]
    
    opt1_percs.append(perc1.text)
    optX_percs.append(percX.text)
    opt2_percs.append(perc2.text)

    if len(stats_info_html) > 3:
        odds1 = stats_info_html[3]
        oddsX = stats_info_html[4]
        odds2 = stats_info_html[5]

        opt1_odds.append(odds1.text)
        optX_odds.append(oddsX.text)
        opt2_odds.append(odds2.text)
    else:
        opt1_odds.append("-")
        optX_odds.append("-")
        opt2_odds.append("-")

powerplay_data = {
    "teams_home_array": teams_home,
    "teams_away_array": teams_away,
    "opt1_percs_array": opt1_percs,
    "optX_percs_array": optX_percs,
    "opt2_percs_array": opt2_percs,
    "opt1_odds_array": opt1_odds,
    "optX_odds_array": optX_odds,
    "opt2_odds_array": opt2_odds,
}



update = True
powerplay_data_fname = 'output/powerplay_data.json'

if os.path.exists(powerplay_data_fname):
    with open(powerplay_data_fname, 'r') as f:
        old_powerplay_data = json.load(f)
        if old_powerplay_data != powerplay_data:
            update = True
else:
    update = True

if update:
    print("Updating powerplay data")
    with open(powerplay_data_fname, 'w') as f:
        f.write(json.dumps(powerplay_data, indent=4))

    rendered_html = p_template.render(powerplay_data)

    with open('html/rendered/powerplay.html', 'w') as f:
        f.write(rendered_html)