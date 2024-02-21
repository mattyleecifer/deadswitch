# pip install google-auth google-auth-oauthlib google-auth-httplib2 google-api-python-client

import base64
from email.mime.text import MIMEText
from google_auth_oauthlib.flow import InstalledAppFlow
from googleapiclient.discovery import build
from requests import HTTPError
import argparse

def send_email(subject, recipients, text):
    SCOPES = [
            "https://www.googleapis.com/auth/gmail.send"
        ]
    flow = InstalledAppFlow.from_client_secrets_file('credentials.json', SCOPES)
    creds = flow.run_local_server(port=0)

    service = build('gmail', 'v1', credentials=creds)
    message = MIMEText(text)
    message['to'] = recipients
    message['subject'] = subject
    create_message = {'raw': base64.urlsafe_b64encode(message.as_bytes()).decode()}

    try:
        message = (service.users().messages().send(userId="me", body=create_message).execute())
        print(F'sent message to {message} Message Id: {message["id"]}')
    except HTTPError as error:
        print(F'An error occurred: {error}')
        message = None

def parse_args():
    # create the argument parser
    parser = argparse.ArgumentParser(description='Send a message')

    # add the arguments
    parser.add_argument('-d', '--days', type=int, required=True,
                        help='number of days before the message is sent')
    parser.add_argument('-a', '--addrecipient', type=str, action='append',
                        help='add a recipient to the message')
    parser.add_argument('-m', '--message', type=str, required=True,
                        help='message to be sent')
    parser.add_argument('-auth', '--auth', type=str, required=True,
                        help='authorization key')

    # parse the arguments
    args = parser.parse_args()

    return args
    # extract the recipients as a list


def main():
    # get flag from args (eg how many days)
    # flag should have owner and recipient emails -owner {} -addrecipient {} - multiple recipients
    args = parse_args()
    days = args.days
    recipients = args.addrecipient
    message = args.message
    auth = args.auth
    # set timer from the arg - emails user to reminder at 1/2 and 3/4 time
    # at timer, send email
    # load webpage with just the auth - every time the auth is valid, timer resets
    return