🎓 LetStudy CLI

A productivity-focused Command Line App built for Computer Science students and developers — helping you manage tasks, notes, reflections, and focus sessions right from your terminal.

🧩 Features
✅ Task & To-Do Management
Create new tasks with optional due dates
Mark tasks as done
Automatic reminder for upcoming due tasks
View active and completed tasks

🧠 Reflection
After marking a task as done, the CLI prompts you to reflect on it
Stores insights, improvements, and a satisfaction rating (1–5)
You can list all reflections or view by specific task ID

📝 Quick Notes
Instantly jot down quick ideas or thoughts
Search, view, and delete notes easily
Timestamped automatically

⏳ Pomodoro Timer
Focus timer with pause/resume control
Visual progress bar and completion sound notification
Session history saved in SQLite

🤖 Fun Commands
dadjoke — Get a random dad joke
quote — Fetch motivational quotes (from ZenQuotes API)

``PREREQUISITES``
Make sure golang installed in your device. Install it through this link ( https://go.dev/doc/install ) . 
After it is installed make sure that you have go path (C:\Users\<username>\go\bin) in the environment path
  

Installation
```bash
go install github.com/bradleyalexa/letstudy@latest
```

Optional
Create a shorter command 'lst'
```bash
git clone https://github.com/bradleyalexa/letstudy.git
cd letstudy
go build -o "$env:USERPROFILE\go\bin\lst.exe"
```

Usage Guide

Todo / Task
```bash
letstudy task new          # Create a new task
letstudy task list         # Show active tasks
letstudy task list done    # Show completed tasks
letstudy task markdone 3   # Mark task as done (includes reflection)
letstudy task remind       # Show upcoming due tasks
```

Reflection
```bash
letstudy reflect list      # View all reflections (task ID + title)
letstudy reflect view 3    # View reflection details for reflection ID 3
```

Quick Notes
```bash
letstudy note add          # Create a new quick note
letstudy note list         # List all notes
letstudy note view 2       # View note by ID
letstudy note delete 2     # Delete a note
letstudy note search "AI"  # Search notes by keyword
```

Pomodoro
```bash
letstudy pomodoro --minutes 25 --type focus
letstudy pomodorohistory   # Show pomodoro session history
```

Fun Commands
```bash
letstudy dadjoke           # Get a random dad joke
letstudy quote             # Get a random motivational quote
```

Example WorkFlow
```bash
letstudy task new
letstudy task list
letstudy task markdone 3
letstudy reflect list
letstudy reflect view 3
letstudy note add
letstudy pomodoro --minutes 25
```

📦 Data Storage

All data is stored locally in:
```
sqlite-database.db
```
Automatically created at first launch.
No external dependencies required.

© 2025 Alexander Bradley
