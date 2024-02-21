# pitt-internship-api

Rest api that scrapes the [Summer 2024 Tech Internships by Pitt CSC & Simplify](https://github.com/SimplifyJobs/Summer2024-Internships) page for internships

## Get job listings

### Request

`GET /`

### Response

    {
        "listings": [
            {
                "company": "Example company",
                "role": "Software Engineer Intern",
                "locations": [
                    "New York, NY"
                ],
                "link": "https://www.example.com",
                "date": "2024-02-20T00:00:00Z"
            },
        ]
    }
