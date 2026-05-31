# Frontend Implementation Plan

## Execution Order

### Phase 1: Project Scaffolding
- [ ] Initialize Vite + React + TypeScript project in `web/`
- [ ] Install dependencies (antd, react-router-dom, zustand, axios, socket.io-client)
- [ ] Configure Vite proxy for API
- [ ] Set up directory structure (pages, components, services, stores, hooks, types, utils)
- [ ] Create Dockerfile for production build

### Phase 2: Foundation Layer
- [ ] Define TypeScript types (`types/`) matching backend models
- [ ] Create Axios instance with interceptors (token injection, error handling)
- [ ] Create API service modules (`services/`) for each backend module
- [ ] Set up Zustand stores (auth, user, notifications)
- [ ] Create route configuration with auth guard
- [ ] Create shared Layout component (header, sidebar, content)

### Phase 3: Auth Module
- [ ] Login page (phone + code input, user type selection)
- [ ] SMS code countdown timer
- [ ] Token persistence (localStorage)
- [ ] Auto-redirect on login state

### Phase 4: Item Module
- [ ] Home page with item grid/list
- [ ] Category navigation
- [ ] Search with filters (keyword, category, condition, price range)
- [ ] Item detail page
- [ ] Item publish form (title, description, images upload, category, condition, want_items)
- [ ] My items management page (edit, toggle status, delete)

### Phase 5: Trade Module
- [ ] Initiate trade dialog (select offered item, set barter coin amount, message)
- [ ] Trade list page (tabs: all/pending/accepted/completed)
- [ ] Trade detail page with action buttons (accept/reject/complete/cancel)

### Phase 6: Message Module
- [ ] Conversation list page
- [ ] Chat interface (message bubbles, input, send)
- [ ] WebSocket connection for real-time messages
- [ ] Unread count badge

### Phase 7: Wallet & Profile
- [ ] Wallet page (balance display, transaction history)
- [ ] Profile edit page
- [ ] Address management (CRUD)
- [ ] User reviews page

### Phase 8: Notifications
- [ ] Notification list page
- [ ] Unread count in header
- [ ] Mark read / mark all read

## Validation Commands

```bash
cd web && npm run build    # TypeScript compilation check
cd web && npm run lint     # ESLint check
```

## Notes

- Use Ant Design's responsive Grid system (Row/Col with breakpoints)
- Use Ant Design's Form component with validation rules
- Image upload: use Ant Design Upload component, preview with Modal
- Mobile-first: test at 375px width as primary viewport
