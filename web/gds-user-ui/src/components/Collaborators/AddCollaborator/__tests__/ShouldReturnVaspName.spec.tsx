/* eslint-disable @typescript-eslint/no-unused-vars */
import React from 'react';
import { waitFor } from '@testing-library/react';
import { useQuery } from '@tanstack/react-query';
import userEvent from '@testing-library/user-event';
import { dynamicActivate } from 'utils/i18nLoaderHelper';
import AddCollaboratorForm from '../AddCollaboratorForm';
import { act, render } from 'utils/test-utils';
import { useSelector } from 'react-redux';
import { userSelector } from 'modules/auth/login/user.slice';
// mock use selector
jest.mock('react-redux', () => ({
  ...jest.requireActual('react-redux'),
  useSelector: jest.fn()
}));

// mock userSelector to return user
jest.mock('modules/auth/login/user.slice', () => ({
  ...jest.requireActual('modules/auth/login/user.slice'),
  userSelector: jest.fn().mockReturnValue({
    user: {
      vasp: {
        id: '1',
        name: 'vasp-test'
      }
    }
  })
}));

function renderComponent() {
  const Props = {
    onCloseModal: jest.fn()
  };

  return render(<AddCollaboratorForm {...Props} />);
}

describe('User Organization', () => {
  it('should return vasp name', () => {
    const { queryByTestId } = renderComponent();
    expect(queryByTestId('vasp-name')).toBeNull();
  });
});
