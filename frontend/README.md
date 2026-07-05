# Scout Frontend

React 19 + TypeScript + Vite frontend for the Scout greenhouse monitoring system.

## Quick Start

### Installation
```bash
npm install --legacy-peer-deps
```

### Development
```bash
npm run dev
```
Server runs at `http://localhost:5173` with API proxy to `http://localhost:8080`

### Build
```bash
npm run build
```
Output: `dist/`

### Testing
```bash
npm test
```

## Architecture

- **React 19** with TypeScript
- **Redux Toolkit** for state management
- **RTK Query** for API data fetching
- **Tailwind CSS** for styling
- **Konva.js** for interactive greenhouse map
- **Vitest** + **React Testing Library** for testing

## Project Structure

See [CLAUDE.md](./CLAUDE.md) for detailed conventions and structure.

## Implementation Phases

- **Phase A** ✅ Scaffold (Vite + dependencies)
- **Phase B** Types & RTK Query setup
- **Phase C** Redux store + filters
- **Phase D-F** Gallery, filters, map components
- **Phase G** Responsive layout
- **Phase H** Component tests
- **Phase I** Docker support
- **Phase J** Documentation & env vars

## API Integration

- Backend: `http://localhost:8080` (proxied via `/api`)
- Database: SQLite (read-only, WAL mode)
- Storage: MinIO S3-compatible
- Endpoints: Photos, thumbnails, health check
