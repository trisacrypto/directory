import type { PersistedState, MigrationManifest } from './types';
export default function createMigrate(migrations: MigrationManifest, config?: {
    debug: boolean;
}): (state: PersistedState, currentVersion: number) => Promise<PersistedState>;
