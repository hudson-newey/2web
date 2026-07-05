import { createAction, createReducer, createSelector, on } from '@ngrx/store';

export type CounterState = {
  readonly count: number;
  readonly nestedObject: {
    readonly nestedCount: number;
    readonly lastUpdated: string;
    readonly depth: {
      readonly intermediate: {
        readonly details: {
          readonly id: string;
          readonly tags: readonly [string, string, number];
        };
      };
    };
  };
  readonly userProfiles: ReadonlyArray<{
    readonly id: string;
    readonly name: string;
    readonly email?: string;
    readonly role: 'admin' | 'editor' | 'viewer';
    readonly preferences: {
      readonly notifications: {
        readonly email: boolean;
        readonly push: boolean;
      };
      readonly theme: 'light' | 'dark' | 'system';
      readonly languages: readonly string[];
    };
    readonly recentActivity: ReadonlyArray<
      | { readonly type: 'login'; readonly timestamp: string }
      | { readonly type: 'edit'; readonly documentId: string; readonly diff: { readonly changed: number; readonly added: number; readonly removed: number } }
      | { readonly type: 'comment'; readonly commentId: string; readonly content: string }
    >;
  }>;
  readonly entityMap: Readonly<Record<string, {
    readonly entityId: string;
    readonly metadata: {
      readonly createdAt: string;
      readonly updatedAt: string;
      readonly attributes: Record<string, string | number | boolean>;
    };
    readonly tags: readonly string[];
  }>>;
  readonly statistics: {
    readonly totalSessions: number;
    readonly recentScores: readonly number[];
    readonly histogram: readonly [number, number, number, number, number];
  };
};

const now = new Date().toISOString();

const initialState: CounterState = {
  count: 0,
  nestedObject: {
    nestedCount: 0,
    lastUpdated: now,
    depth: {
      intermediate: {
        details: {
          id: 'detail-001',
          tags: ['alpha', 'beta', 42],
        },
      },
    },
  },
  userProfiles: [
    {
      id: 'u1',
      name: 'Avery',
      email: 'avery@example.com',
      role: 'admin',
      preferences: {
        notifications: { email: true, push: false },
        theme: 'dark',
        languages: ['en', 'fr'],
      },
      recentActivity: [
        { type: 'login', timestamp: now },
        { type: 'comment', commentId: 'c123', content: 'Reviewed the new design.' },
      ],
    },
    {
      id: 'u2',
      name: 'Riley',
      role: 'editor',
      preferences: {
        notifications: { email: false, push: true },
        theme: 'system',
        languages: ['es'],
      },
      recentActivity: [
        { type: 'edit', documentId: 'doc-456', diff: { changed: 12, added: 3, removed: 1 } },
      ],
    },
  ],
  entityMap: {
    'e-1': {
      entityId: 'e-1',
      metadata: {
        createdAt: now,
        updatedAt: now,
        attributes: { isActive: true, score: 97, label: 'primary' },
      },
      tags: ['alpha', 'customer'],
    },
    'e-2': {
      entityId: 'e-2',
      metadata: {
        createdAt: now,
        updatedAt: now,
        attributes: { isActive: false, score: 12, label: 'secondary' },
      },
      tags: ['beta', 'internal'],
    },
  },
  statistics: {
    totalSessions: 13,
    recentScores: [88, 92, 74, 99],
    histogram: [1, 4, 5, 2, 1],
  },
};

export const increment = createAction('[Counter] Increment');

export const counterReducer = createReducer(
  initialState,
  on(increment, (state) => ({
    ...state,
    count: state.count + 1,
    nestedObject: {
      ...state.nestedObject,
      nestedCount: state.nestedObject.nestedCount + Math.random(),
      lastUpdated: new Date().toISOString(),
    },
    statistics: {
      ...state.statistics,
      totalSessions: state.statistics.totalSessions + 1,
      recentScores: [...state.statistics.recentScores, Math.round(Math.random() * 100)],
    },
  })),
);

export const counterState = createSelector(
  (state: CounterState) => state.count,
  (count) => count,
);
