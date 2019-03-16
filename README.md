# colly-draft-prospects
Source code for web scraping football draft prospects from the DraftTek NFL Big Board.

The scraper is written in Golang and uses the [Colly](http://go-colly.org/) scraper. The binary file in the repo is compiled for Linux, but it could be compiled to use in a different operating system if needed.

Once the ranks have been scraped, I use [csvkit](https://csvkit.readthedocs.io]) to merge all of the files and join them together with information about the locations of the schools.  The csvkit commands are in the csvkitcommands.txt files.

Once the ranks have been assembled, I use [OpenRefine](http://openrefine.org/) to clean the data for consistency.  The data cleaning steps are contained in openRefineDataMerge.json.
