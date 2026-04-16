# Python Service Technical Documentation

## Overview

The Python service is a FastAPI application that provides AI-powered resume generation, cover letter creation, and job description parsing.

## Directory Structure

```
python-service/
├── main.py                         # FastAPI application
├── src/
│   ├── resume_generator.py         # Resume generation logic
│   ├── cover_letter_generator.py   # Cover letter generation
│   └── job_parser.py               # Job description parsing
├── data_folder/                    # Output directory
│   └── output/                     # Generated PDFs
├── requirements.txt                # Python dependencies
└── Dockerfile                      # Container configuration
```

## Main Application (`main.py`)

### FastAPI Setup

```python
app = FastAPI(title="Job Applier Resume Service", version="1.0.0")

app.add_middleware(
    CORSMiddleware,
    allow_origins=["*"],
    allow_credentials=True,
    allow_methods=["*"],
    allow_headers=["*"],
)
```

### Endpoints

| Endpoint | Method | Description |
|----------|--------|-------------|
| `/health` | GET | Health check |
| `/generate-resume` | POST | Generate resume PDF |
| `/generate-cover-letter` | POST | Generate cover letter PDF |
| `/parse-job` | POST | Parse job description |
| `/styles` | GET | List available styles |
| `/templates` | GET | List available templates |

## Resume Generator (`src/resume_generator.py`)

### Class Structure

```python
class ResumeGenerator:
    def __init__(self):
        self.llm_api_key = os.getenv("LLM_API_KEY", "")
        self.llm_model = os.getenv("LLM_MODEL", "gpt-4o-mini")
    
    def generate(
        self,
        resume_yaml: str,
        style: str = "modern",
        job_description: Optional[str] = None,
        output_path: str = "output/resume.pdf"
    ) -> Dict[str, Any]:
        # 1. Parse YAML
        # 2. Tailor to job (if provided)
        # 3. Generate HTML
        # 4. Save HTML
        # 5. Generate PDF
        # 6. Return result
```

### Resume Styles

| Style | Description |
|-------|-------------|
| `modern` | Clean, contemporary design with accent colors |
| `classic` | Traditional, professional layout |
| `minimal` | Simple, elegant with focus on content |
| `creative` | Bold design for creative roles |
| `professional` | Corporate-focused, ATS-friendly |

### HTML Generation

Generates HTML from structured resume data:

```python
def _generate_html(self, resume_data: Dict, style: str) -> str:
    # Renders sections:
    # - Header (name, contact, links)
    # - Summary
    # - Experience
    # - Education
    # - Skills
    # - Projects
```

### LLM Integration

When `LLM_API_KEY` is set, uses LangChain to tailor resumes:

```python
from langchain_openai import ChatOpenAI
from langchain.prompts import ChatPromptTemplate

llm = ChatOpenAI(
    model=self.llm_model,
    api_key=self.llm_api_key,
    temperature=0.7
)
```

### PDF Generation

Two methods:
1. **Selenium + Chrome**: Headless Chrome prints HTML to PDF
2. **ReportLab fallback**: Basic PDF if Chrome unavailable

```python
def _generate_pdf(self, html_content: str, output_path: str):
    try:
        # Try Selenium
        driver = webdriver.Chrome(options=chrome_options)
        driver.print_page(output_path)
    except Exception:
        # Fallback to ReportLab
        doc = SimpleDocTemplate(output_path, pagesize=letter)
```

## Cover Letter Generator (`src/cover_letter_generator.py`)

### Class Structure

```python
class CoverLetterGenerator:
    def generate(
        self,
        resume_text: str,
        job_description: str,
        job_url: Optional[str] = None,
        company_name: Optional[str] = None,
        job_title: Optional[str] = None,
        output_path: str = "output/cover_letter.pdf"
    ) -> Dict[str, Any]:
        # Generate content
        # Wrap in HTML
        # Save files
        # Return result
```

### Generation Methods

1. **LLM-powered**: When API key available
   - Uses LangChain with ChatOpenAI
   - Generates personalized content
   - Tailored to job description

2. **Template-based**: Fallback
   - Basic template with placeholders
   - Extracts highlights from resume

### LLM Prompt

```python
prompt = ChatPromptTemplate.from_template(
    "Write a professional cover letter for:\n\n"
    "Job Title: {job_title}\n"
    "Company: {company_name}\n"
    "Job Description: {job_description}\n\n"
    "Candidate Resume: {resume}\n\n"
    "Requirements:\n"
    "1. Opens with enthusiasm\n"
    "2. Highlights relevant experience\n"
    "3. Shows understanding of requirements\n"
    "4. Closes with call to action"
)
```

## Job Parser (`src/job_parser.py`)

### Class Structure

```python
class JobParser:
    def parse_url(self, url: str) -> Dict[str, Any]:
        if self.llm_api_key:
            return self._parse_with_llm(url)
        else:
            return self._parse_basic(url)
```

### Parsing Methods

1. **LLM-powered** (`_parse_with_llm`):
   - Uses Selenium to fetch page content
   - Sends to LLM for structured extraction
   - Returns JSON with fields

2. **Basic** (`_parse_basic`):
   - Uses Selenium to fetch page
   - Regex-based extraction
   - Pattern matching for common fields

### Extracted Fields

```python
{
    "title": str,
    "company": str,
    "location": str,
    "description": str,
    "requirements": List[str],
    "responsibilities": List[str],
    "salary": Optional[str],
    "remote": bool
}
```

## Data Flow

### Resume Generation Flow

```
Request → Parse YAML → LLM Tailoring → HTML Generation → PDF Generation → Response
```

1. Parse YAML resume data
2. If job description provided, tailor content via LLM
3. Generate HTML with selected style CSS
4. Save HTML to disk
5. Convert HTML to PDF (Selenium or ReportLab)
6. Return PDF path and metadata

### Cover Letter Flow

```
Request → LLM Generation → HTML Wrapping → PDF Generation → Response
```

1. Generate cover letter content via LLM
2. Wrap in styled HTML template
3. Convert to PDF
4. Return path and metadata

## Environment Variables

| Variable | Default | Description |
|----------|---------|-------------|
| `LLM_API_KEY` | - | OpenAI API key |
| `LLM_MODEL` | `gpt-4o-mini` | Model name |
| `CHROME_BIN` | - | Chrome binary path |
| `CHROMEDRIVER_PATH` | - | ChromeDriver path |

## Dependencies

```txt
fastapi==0.115.0          # Web framework
uvicorn[standard]==0.30.0 # ASGI server
pydantic==2.9.0           # Data validation
pyyaml==6.0.2             # YAML parsing
reportlab==4.2.0          # PDF generation
selenium==4.25.0          # Browser automation
webdriver-manager==4.0.2  # ChromeDriver management
langchain==0.2.11         # LLM framework
langchain-openai==0.1.20  # OpenAI integration
httpx==0.27.0             # HTTP client
python-dotenv==1.0.1      # Environment variables
```

## Error Handling

- Pydantic validation errors return 400
- Generation failures return 500 with error message
- Graceful fallbacks when LLM unavailable
- Logging via Python's logging module

## Chrome/Selenium Setup

### Docker

```dockerfile
RUN apt-get update && apt-get install -y \
    chromium \
    chromium-driver \
    fonts-liberation

ENV CHROME_BIN=/usr/bin/chromium
ENV CHROMEDRIVER_PATH=/usr/bin/chromedriver
```

### Local Development

Install Chrome and ChromeDriver, or use webdriver-manager:
```python
from webdriver_manager.chrome import ChromeDriverManager
driver = webdriver.Chrome(ChromeDriverManager().install())
```

## Testing

```bash
# Run tests
python -m pytest tests/ -v

# Run with coverage
python -m pytest tests/ --cov=src
```

## Performance Considerations

- Chrome startup is slow (~2s)
- Consider caching generated PDFs
- LLM calls add latency (~1-3s)
- Async FastAPI handlers for concurrency

## Production Deployment

```bash
# With multiple workers
uvicorn main:app --host 0.0.0.0 --port 8001 --workers 4

# With Gunicorn
gunicorn main:app -w 4 -k uvicorn.workers.UvicornWorker
```
