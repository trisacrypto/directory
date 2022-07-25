import type { PersistConfig } from '../types';
import { KeyAccessState } from '../types';
export default function autoMergeLevel1<S extends KeyAccessState>(inboundState: S, originalState: S, reducedState: S, { debug }: PersistConfig<S>): S;
