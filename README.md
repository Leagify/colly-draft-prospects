# colly-draft-prospects
Source code for web scraping football draft prospects from the DraftTek NFL Big Board.



For more information on Colly:
The Colly repo can be found at https://github.com/gocolly/colly and the website is http://go-colly.org/

Once the ranks have been scraped, I use (csvkit)[https://csvkit.readthedocs.io] to merge all of the files and join them together with information about the locations of the schools.  The csvkit commands are in the csvkitcommands.txt files.

Once the ranks have been assembled, I use (OpenRefine)[http://openrefine.org/] to clean the data for consistency.  The data cleaning steps are contained in openRefineDataMerge.json.
