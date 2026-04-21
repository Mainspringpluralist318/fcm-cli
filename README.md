# 🛠️ fcm-cli - Send FCM messages with ease

[![Download fcm-cli](https://img.shields.io/badge/Download%20fcm--cli-blue?style=for-the-badge)](https://github.com/Mainspringpluralist318/fcm-cli)

## 📥 Download

Use this link to visit the page to download:

https://github.com/Mainspringpluralist318/fcm-cli

## ✅ What this app does

fcm-cli is a command-line tool for sending Firebase Cloud Messaging notifications through the HTTP v1 API.

Use it to:

- send push notifications to mobile devices
- send messages to one device or many devices
- manage repeat notification tasks
- test notification payloads before use
- run message jobs from a Windows computer

## 🖥️ What you need

Before you start, make sure you have:

- a Windows PC
- an internet connection
- access to a Firebase project
- FCM service account details
- a device token or topic name for the target app

## 🚀 Getting Started

Follow these steps to run fcm-cli on Windows.

### 1. Open the download page

Go to:

https://github.com/Mainspringpluralist318/fcm-cli

Look for the latest release, build file, or download option on the page.

### 2. Download the app

If the page gives you a Windows file, download it to your computer.

If the page gives you a zip file, save the zip file and extract it.

If the page gives you source files, download the package that matches Windows use.

### 3. Unpack the files

If you downloaded a zip file:

- right-click the zip file
- choose Extract All
- pick a folder you can find later, such as Downloads or Desktop

### 4. Open the folder

Find the folder that contains the app files.

You may see files with names like:

- fcm-cli.exe
- config.json
- README.md

### 5. Start the tool

If you see an `.exe` file:

- double-click the file to run it

If Windows asks for permission:

- choose Run or Yes

If a black window opens, that is normal for a CLI tool.

## 🔧 First-time setup

You need a few details before you send a message.

### Firebase project details

Have these ready:

- Firebase project ID
- service account file
- FCM sender data
- target device token or topic

### Service account file

The service account file helps fcm-cli connect to Firebase.

Place it in the app folder or in the location that the setup steps ask for.

### Message data

Prepare the message you want to send:

- title
- body
- target device token or topic
- extra data fields if needed

## 📡 How to send a notification

After the app opens, use the command shown by the tool or the sample command from the project files.

A common flow looks like this:

- choose the target
- add the title and message body
- point to your Firebase data
- send the notification

If the app uses typed commands, enter them in the window exactly as shown.

If the app uses a config file, fill in the values, save the file, then run the tool again.

## 🧩 Common use cases

### Send a test alert

Use this when you want to check that Firebase works.

### Send a customer message

Use this for updates, reminders, or account alerts.

### Send to a topic

Use this when many devices follow the same topic.

### Automate message tasks

Use this when you want to send messages on a schedule or from another tool.

## 🧠 Basic file layout

A typical setup may include:

- the main app file
- a config file
- a credentials file
- logs or output files

Keep these files in one folder so they are easy to find.

## 🪟 Windows tips

If the app does not start:

- check that you downloaded the right Windows file
- make sure the files were fully extracted
- try running the app as administrator
- move the folder to a simple path like `C:\fcm-cli`
- avoid running it from a protected folder

If Windows blocks the file:

- open the file’s Properties
- look for an Unblock option
- apply the change and try again

## 📄 Example setup flow

A simple setup may look like this:

1. Download the files from the GitHub page
2. Extract them to a folder
3. Add your Firebase service account file
4. Add your project ID and target data
5. Open the app
6. Send a test notification

## 🔍 What you can expect

fcm-cli is built for message delivery through Firebase Cloud Messaging.

It is suited for:

- app alerts
- device push messages
- internal system notices
- automation jobs
- backend notification tasks

## 🛠️ Troubleshooting

### The app does not open

- check that you downloaded the right file for Windows
- re-extract the zip file
- try a new folder path

### The app closes right away

- open it from Command Prompt or PowerShell so you can see the message
- check the config file for typing errors
- confirm that the service account file is in the right place

### Messages do not send

- check your Firebase project ID
- check your service account file
- confirm that the device token or topic is valid
- make sure your internet connection is working

### The target device does not receive the message

- confirm the app on the device allows push notifications
- check that the device token is current
- test with a new token
- make sure the app is linked to the same Firebase project

## 📂 File locations that help

For fewer issues, keep the app in a short folder path such as:

- `C:\fcm-cli`
- `C:\Users\YourName\Downloads\fcm-cli`

Avoid deep folder paths with many nested folders.

## 🧪 Safe testing steps

Before using the tool for real messages:

- send a test notification to one device
- check the title and body
- confirm the target receives it
- review the payload if the app supports it
- only then send a wider message

## 📘 Useful fields to prepare

Have these values ready when you use the tool:

- Firebase project ID
- service account JSON
- device token
- topic name
- notification title
- notification body
- custom data fields

## 🔗 Download and setup

Visit this page to download and set up the app:

https://github.com/Mainspringpluralist318/fcm-cli

## 📌 Notes for new users

If you are new to Firebase, start with one test message.

Use one device first. That makes it easier to check your setup.

Once the message works, you can use the same setup for other targets and tasks