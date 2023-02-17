import { ErrorBoundary } from 'react-error-boundary';

import TimeAgo from '@/components/TimeAgo';
import { render, screen, waitFor } from '@/utils/test-utils';

const FAKE_SYSTEM_TIMESTAMP = 1638305340000; // 2021-11-30T20:49:00.000Z

describe('<TimeAgo />', () => {
  beforeAll(() => {
    jest.useFakeTimers('modern').setSystemTime(FAKE_SYSTEM_TIMESTAMP);
  });

  afterAll(() => {
    jest.useRealTimers();
  });

  it('should increase time each 1second', async () => {
    render(<TimeAgo time={FAKE_SYSTEM_TIMESTAMP} />);
    const timeElement = screen.getByTestId(/time/i);
    expect(timeElement.textContent).toBe('a few seconds ago');
  });

  it('should throw error when null is passed', async () => {
    const fallbackRender = jest.fn(() => null);

    render(
      <ErrorBoundary fallbackRender={fallbackRender}>
        <TimeAgo time={null} />
      </ErrorBoundary>
    );

    await waitFor(() => expect(fallbackRender).toHaveBeenCalled());
  });

  it('should throw error when undefined is passed', async () => {
    const fallbackRender = jest.fn(() => null);

    render(
      <ErrorBoundary fallbackRender={fallbackRender}>
        <TimeAgo time={undefined} />
      </ErrorBoundary>
    );

    await waitFor(() => expect(fallbackRender).toHaveBeenCalled());
  });
});
