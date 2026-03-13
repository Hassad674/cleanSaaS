# CleanSaaS Mobile

Flutter mobile app for CleanSaaS. Connects to the same Go backend as the web frontend.

## Prerequisites

- Flutter SDK >= 3.16.0
- Dart SDK >= 3.2.0
- Android Studio or Xcode (for emulators)

## Setup

### Option A: Fresh project (recommended)

```bash
# Create Flutter project
flutter create --org com.cleansaas --project-name cleansaas_mobile temp_project

# Copy our lib/ into the created project
cp -r lib/ temp_project/lib/
cp pubspec.yaml temp_project/
cp analysis_options.yaml temp_project/

# Replace the temp project
mv temp_project/* .
rm -rf temp_project

# Install dependencies
flutter pub get
```

### Option B: Use this directory directly

If you already have Flutter configured:

```bash
flutter pub get
flutter run
```

## Configuration

Edit `lib/config/constants.dart` to set your API URL:

- **Android emulator**: `http://10.0.2.2:8081` (default)
- **iOS simulator**: `http://localhost:8081`
- **Physical device**: Use your machine's local IP (e.g., `http://192.168.1.X:8081`)
- **Production**: Your deployed backend URL

## Architecture

Feature-based architecture matching the web frontend:

```
lib/
├── main.dart              → Entry point
├── app.dart               → MaterialApp with router + theme
├── config/                → Theme, constants, router
├── core/                  → Shared infrastructure
│   ├── api/               → HTTP client (Dio) with auth interceptor
│   ├── storage/           → Secure token storage
│   └── widgets/           → Reusable UI components
└── features/              → Business logic modules
    └── auth/              → Authentication feature
        ├── models/        → Data models
        ├── providers/     → Riverpod state management
        ├── repositories/  → API calls
        └── screens/       → UI screens
```

### Adding a new feature

1. Create `features/<name>/models/` with data models
2. Create `features/<name>/repositories/` with API calls
3. Create `features/<name>/providers/` with Riverpod state
4. Create `features/<name>/screens/` with UI
5. Add routes in `config/router.dart`

Features never import from other features. Shared code lives in `core/`.

## Running

```bash
# Development
flutter run

# Run on specific device
flutter run -d chrome        # Web (for quick testing)
flutter run -d emulator-5554 # Android emulator
flutter run -d iPhone        # iOS simulator

# Build
flutter build apk            # Android APK
flutter build ios             # iOS (requires macOS + Xcode)
```

## Testing

```bash
flutter test
```
