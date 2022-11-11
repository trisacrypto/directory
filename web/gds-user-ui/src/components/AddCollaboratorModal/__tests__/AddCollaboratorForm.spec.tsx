/* eslint-disable @typescript-eslint/no-unused-vars */
import React from 'react';
import { waitFor } from '@testing-library/react';
import { QueryClient, QueryClientProvider, useQuery } from '@tanstack/react-query';
import userEvent from '@testing-library/user-event';
import { dynamicActivate } from 'utils/i18nLoaderHelper';
import AddCollaboratorForm from '../AddCollaboratorForm';
import { act, render } from 'utils/test-utils';
export function useCustomHook() {
  return useQuery({ queryKey: ['customHook'], queryFn: () => 'Hello' });
}

// export function useMutationHook() {
//   return useMutation({ mutationKey: ['customHook'], mutationFn: () => 'Hello' });
// }

const queryClient = new QueryClient();

function renderComponent() {
  const Props = {
    onCloseModal: jest.fn()
  };

  return render(
    <QueryClientProvider client={queryClient}>
      <AddCollaboratorForm {...Props} />
    </QueryClientProvider>
  );
}

describe('AddCollaboratorForm', () => {
  beforeAll(() => {
    act(() => {
      dynamicActivate('en');
    });
  });
  it('should render', () => {
    const { container } = renderComponent();
    expect(container).toMatchSnapshot();
  });
});
