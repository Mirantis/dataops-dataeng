import os
from atlassian import Jira

import json
with open('/dataeng/.secrets/secrets.json','r') as f:
      config = json.load(f)

jira_instance = Jira(
    url = "https://mirantis.jira.com",
    username = (config['user']['username']),
    password = (config['user']['password'])
)

projects = jira_instance.get_all_projects(included_archived=None)
print(projects)
