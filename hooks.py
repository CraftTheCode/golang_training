import os
import re
from pathlib import Path

# Silence upstream Material for MkDocs 2.0 announcement banner
os.environ["NO_MKDOCS_2_WARNING"] = "1"

from mkdocs.structure.files import File

def on_files(files, config):
    """
    Dynamically registers all repository markdown files and generates virtual
    code viewer pages for Go, SQL, Docker, and Makefiles with syntax highlighting.
    """
    repo_root = os.path.abspath(".")
    for p in Path(".").glob("**/*"):
        if p.is_dir() or any(part.startswith(".") or part in ("11-practice", "site", "__pycache__", "docs") for part in p.parts):
            continue
            
        rel_posix = p.as_posix()
        ext = p.suffix.lower()
        name = p.name.lower()
        
        if ext == ".md":
            files.append(File(
                path=rel_posix,
                src_dir=repo_root,
                dest_dir=config["site_dir"],
                use_directory_urls=config["use_directory_urls"]
            ))
        else:
            lang = None
            if ext == ".go":
                lang = "go"
            elif ext == ".sql":
                lang = "sql"
            elif name == "dockerfile":
                lang = "dockerfile"
            elif name in ("docker-compose.yml", "docker-compose.yaml"):
                lang = "yaml"
            elif name == "makefile":
                lang = "makefile"
                
            if lang:
                try:
                    code = p.read_text(encoding="utf-8")
                except Exception:
                    continue
                    
                clean_name = p.name.replace(".", "_")
                if "." not in p.name:
                    clean_name = f"{p.name}_view"
                virt_path = p.parent / f"{clean_name}.md"
                src_uri = str(virt_path).replace("\\", "/")
                
                content = f"# `{p.name}`\n\n"
                content += f"📁 **Source Path**: `{rel_posix}`\n\n"
                content += f"```{lang} linenums=\"1\"\n{code}\n```\n"
                
                files.append(File.generated(config, src_uri, content=content))
                
    return files

def on_page_markdown(markdown, page, config, files):
    """
    Automatically resolves relative folder links (e.g., './01-philosophy-and-runtime/')
    to their target markdown documentation ('./01-philosophy-and-runtime/README.md').
    """
    def replace_folder_link(match):
        folder = match.group(1)
        return f"(./{folder}/README.md)"

    return re.sub(r'\(\./([0-9]{2}-[a-zA-Z0-9_-]+)/\)', replace_folder_link, markdown)
