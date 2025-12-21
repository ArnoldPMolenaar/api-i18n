# 🌐 API I18n: Microservice to Manage Translations

[![Go Version](https://img.shields.io/badge/Go-1.23+-00ADD8?logo=go)](https://go.dev/)
[![Fiber](https://img.shields.io/badge/Fiber-v2-2C8EBB?logo=go)](https://github.com/gofiber/fiber)
[![GORM](https://img.shields.io/badge/GORM-v1-7F52FF?logo=go)](https://gorm.io)
[![Phone Libs](https://img.shields.io/badge/Phone%20Libs-libphonenumber%20%2F%20nyaruka-4CAF50)](https://github.com/nyaruka/phonenumbers)

A lightweight Go microservice that centralizes internationalization (i18n) for applications. It manages apps, categories, keys, and translations, and provides locale/territory lookups and phone utilities. Built with Fiber, GORM, and Valkey/Redis for caching.

---

## 📦 Clone with Submodules

# Developer clone:
```shell
git clone --recurse-submodules <repository_url>
```
# Or if already cloned:
```shell
git submodule update --remote --merge
```

The project uses Git submodules (e.g., shared middleware/utilities). Make sure to clone with `--recurse-submodules` or update submodules after cloning.

---

## 🌱 Initial Seed (CLDR JSON)

This project includes the `cldr-json` submodule under `src/database/fixtures/cldr-json/`. On the first run, the service reads this repository and seeds core locale, territory, and related data into the database.

- First-time seeding duration: typically 15–30 minutes depending on your hardware.
- Automatic skip: if the database already contains the seed data, the process is skipped on subsequent runs.

Tip: keep the submodule updated to get the latest CLDR data.

---

## 🚀 Running with Docker Compose

Build and run the development stack:

```shell
docker compose build dev
docker compose up dev
```

For production, use the `prod` service:

```shell
docker compose build prod
docker compose up -d prod
```

- dev: mounts source, enables hot-reload friendly settings.
- prod: optimized, detached container suitable for deployment.

---

## 🔐 API Endpoints Overview

Private routes require machine authentication (middleware-protected). Public routes are open.

### Private Routes (Machine-Protected)
Base: `/v1`

- Apps
  - `POST /v1/apps/` — Create an app
  - `GET /v1/apps/:name/locales` — Get locales configured for an app
  - `PUT /v1/apps/:name/locales` — Set locales for an app

- Categories
  - `GET /v1/categories/` — List categories
  - `POST /v1/categories/` — Create category
  - `GET /v1/categories/lookup` — Lookup categories
  - `GET /v1/categories/:id` — Get category by ID
  - `PUT /v1/categories/:id` — Update category by ID
  - `DELETE /v1/categories/:id` — Soft-delete category by ID
  - `PUT /v1/categories/:id/restore` — Restore soft-deleted category

- Keys
  - `GET /v1/keys/` — List keys
  - `POST /v1/keys/` — Create key
  - `GET /v1/keys/:id` — Get key by ID
  - `PUT /v1/keys/:id` — Update key by ID
  - `DELETE /v1/keys/:id` — Soft-delete key by ID
  - `PUT /v1/keys/:id/restore` — Restore soft-deleted key

### Public Routes
Base: `/v1`

- Territories
  - `GET /v1/territories/lookup` — Lookup territories (region/country codes)

- Locales
  - `GET /v1/locales/lookup` — Lookup locales (language/script/region combinations)

- Translations
  - `GET /v1/translations/:localeId` — Get translations for a locale

- Phones
  - `GET /v1/phones/lookup` — Phone country codes lookup
  - `GET /v1/phones/validate` — Validate phone number
  - `GET /v1/phones/format` — Format phone number

---

## 🤝 Contributing
We welcome contributions! Please fork the repository and submit a pull request.

## 📝 License
This project is licensed under the MIT License.

## 📞 Contact
For any questions or support, please contact [arnold.molenaar@webmi.nl](mailto:arnold.molenaar@webmi.nl).

<hr />

Made with ❤️ by Arnold Molenaar
