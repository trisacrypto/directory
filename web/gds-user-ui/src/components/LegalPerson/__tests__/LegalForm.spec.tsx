import React from 'react';
import { waitFor } from '@testing-library/react';
import { useQuery } from '@tanstack/react-query';
import userEvent from '@testing-library/user-event';
import { act, render } from 'utils/test-utils';
import { dynamicActivate } from 'utils/i18nLoaderHelper';
import LegalForm from '../LegalForm';
function renderComponent() {
  return render(<LegalForm />);
}

jest.mock('hooks/useFetchCertificateStep', () => ({
  useFetchCertificateStep: () => ({
    certificateStep: {
      form: jest.fn(),
      errors: jest.fn()
    },
    isFetchingCertificateStep: false
  })
}));

describe('LegalForm', () => {
  beforeAll(() => {
    act(() => {
      dynamicActivate('en');
    });
  });

  it('should render', () => {
    const { container } = renderComponent();
    expect(container).toMatchSnapshot();
  });

  it('should render NameIdentifier', () => {
    const { getByTestId } = renderComponent();
    const nameIdentifiers = getByTestId('name-identifier');

    expect(nameIdentifiers).toBeInTheDocument();
  });

  it('should render Address', () => {
    const { getByTestId } = renderComponent();
    const address = getByTestId('legal-adress');

    expect(address).toBeInTheDocument();
  });

  it('should render CountryOfRegistration', () => {
    const { getByTestId } = renderComponent();
    const countryOfRegistration = getByTestId('legal-country-of-registration');

    expect(countryOfRegistration).toBeInTheDocument();
  });

  it('should render NationalIdentification', () => {
    const { getByTestId } = renderComponent();
    const nationalIdentification = getByTestId('legal-name-identification');

    expect(nationalIdentification).toBeInTheDocument();
  });

  it('should render StepButtons', () => {
    const { getByTestId } = renderComponent();
    const stepButtons = getByTestId('step-buttons');

    expect(stepButtons).toBeInTheDocument();
  });
});
