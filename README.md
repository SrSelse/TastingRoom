<a href="https://www.buymeacoffee.com/beanlink" target="_blank"><img src="https://cdn.buymeacoffee.com/buttons/v2/default-yellow.png" alt="Buy Me A Coffee" style="height: 40px !important;width: 145px !important;" ></a>

# TastingRoom


A modern web application for organizing and participating in virtual beer tasting events.

TastingRoom is a webpage that can be used if you want to create a tasting of any kind. I developed it to handle the beer tastings that my brothers and I have every year, but I wanted to extend it from the simple webpage that I made originally. That evolved into what is now TastingRoom, where you can create a user, create a room, invite your friends and let everyone do the ratings on their own device. This way, you can even do the tasting over a zoom call, yet still have a central place to see what everyone thinks about whatever you are tasting.

## 🍺 Features

### For Participants
- **Join tasting rooms** using invitation codes
- **Add beverages** with names, styles, and images
- **Rate beverages** with optional tasting notes
- **View details** including style and images
- **Real-time updates** when new ratings are submitted
- **Responsive design** that works on all devices
- **View average ratings** across all participants

### For Room Admins
- **Create tasting rooms** with names, descriptions, and scheduled dates
- **Manage participants** - add/remove users and assign admin privileges
- **Publish ratings** to make them visible to all participants

## 🚀 Getting Started

### Prerequisites

### Demo

The live webpage is available at [https://tastingroom.online](https://tastingroom.online).

## 📖 Usage

### Creating a Tasting Room

1. Click "Create New Room" from the dashboard
2. Fill in the room details:
   - Room name (required)
   - Description (optional)
   - Planned date (optional)
3. Share the invitation code with participants

### Joining a Tasting Room

1. Get an invitation code from the room admin
2. Enter the code in the invitation field on the dashboard
3. Click "Join Room"

### Adding Beverages

1. Navigate to your tasting room
2. Click "Add Beverage"
3. Fill in beverage details:
   - Name (required)
   - Style (required)
   - Picture URL (optional)
4. Participants can now rate the beverage

### Rating Beverages

1. Select a beverage from the room
2. Enter your rating (0-5)
3. Add optional tasting notes
4. Submit your rating
5. When the admin decides to move on to the next item, follow automatically

### Managing Admins

1. Click "Manage Admins" in your tasting room
2. Use the checkboxes to assign/revoke admin privileges
3. Use the delete button to remove participants

## 🛠 Built With

- **Frontend**: Vue.js 3
- **Styling**: Tailwind CSS
- **Real-time updates**: Centrifugo
- **Authentication**: JWT
- **API**: Custom RESTful API written in Golang

## 📧 Contact

Feedback, improvements and issues are handled in the Issues-tab here on Github.
