/* eslint-disable @typescript-eslint/no-unused-vars */
import React from 'react';
import { waitFor } from '@testing-library/react';
import { useQuery } from '@tanstack/react-query';
import userEvent from '@testing-library/user-event';
import { dynamicActivate } from 'utils/i18nLoaderHelper';
import AddCollaboratorForm from '../AddCollaboratorForm';
import { act, render } from 'utils/test-utils';

function renderComponent() {
  const Props = {
    onCloseModal: jest.fn()
  };

  return render(<AddCollaboratorForm {...Props} />);
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
