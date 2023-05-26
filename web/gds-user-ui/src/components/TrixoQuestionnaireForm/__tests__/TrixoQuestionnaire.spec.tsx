import React from 'react';
import { waitFor, within, act } from '@testing-library/react';

import userEvent from '@testing-library/user-event';
import { render } from 'utils/test-utils';
import { dynamicActivate } from 'utils/i18nLoaderHelper';
import TrixoQuestionnaireForm from '../';

function renderComponent() {
  return render(<TrixoQuestionnaireForm />);
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

describe('TrixoQuestionnaire', () => {
  beforeAll(() => {
    act(() => {
      dynamicActivate('en');
    });
  });

  it('should render', () => {
    const { container } = renderComponent();
    expect(container).toMatchSnapshot();
  });

  it('should render the form', () => {
    const { getByTestId } = renderComponent();
    const trixoQuestionnaire = getByTestId('trixo-form');
    expect(trixoQuestionnaire).toBeInTheDocument();
  });

  it('should render the primary national jurisdiction', () => {
    const { getByText } = renderComponent();
    const primaryNationalJurisdiction = getByText('Primary National Jurisdiction');
    expect(primaryNationalJurisdiction).toBeInTheDocument();
  });

  it('should render the primary national jurisdiction', () => {
    const { getByText } = renderComponent();
    const primaryNationalJurisdiction = getByText('Primary National Jurisdiction');
    expect(primaryNationalJurisdiction).toBeInTheDocument();
  });

  it('should render add jurisdiction button', () => {
    const { getByText } = renderComponent();
    const addJurisdiction = getByText('Add Jurisdiction');
    expect(addJurisdiction).toBeInTheDocument();
  });

  it('should add a new fields when clicking on add jurisdiction button', async () => {
    const { getByTestId } = renderComponent();
    const addJurisdiction = getByTestId('trixo-add-jurisdictions-btn');
    userEvent.click(addJurisdiction);
    await waitFor(() => {
      expect(getByTestId('trixo-country')).toBeInTheDocument();
      expect(getByTestId('trixo-regulator-name')).toBeInTheDocument();
      expect(getByTestId('trixo-add-jurisdictions-btn')).toBeInTheDocument();
    });
  });

  it('should render the required financial transfers question', () => {
    // check if the input after the data-testid is rendered
    const { getByTestId } = renderComponent();
    const wrapper = getByTestId('is-required-financial-transfers');
    const select = within(wrapper).getByRole('combobox');
    expect(select).toBeInTheDocument();
  });

  it('should render the required regulatory program questionnaire', () => {
    const { getByTestId } = renderComponent();
    const wrapper = getByTestId('trixo-rule-policies');
    const select = within(wrapper).getByRole('combobox');
    expect(select).toBeInTheDocument();
  });

  it('should render Conducts KYC radio button', () => {
    const { getByTestId, getByText } = renderComponent();

    const questionnaire = getByText('Conducts KYC before virtual asset transfers');

    expect(questionnaire).toBeInTheDocument();

    const wrapper = getByTestId('trixo-kyc-before-virtual-asset-transfers');
    const radioButtons = within(wrapper).getByRole('checkbox');
    expect(radioButtons).toBeInTheDocument();
    // if we click on the radio button, it should be checked
    userEvent.click(radioButtons);
    expect(radioButtons).toBeChecked();
    // if we click again, it should be unchecked
    userEvent.click(radioButtons);
    expect(radioButtons).not.toBeChecked();
  });

  it('should render Travel Rule radio button', async () => {
    const { getByTestId, getByText } = renderComponent();

    const questionnaire = getByText('Must comply with Travel Rule');

    expect(questionnaire).toBeInTheDocument();

    const wrapper = getByTestId('trixo-must-comply-travel-rule');

    const radioButton = within(wrapper).getByRole('checkbox');

    expect(radioButton).toBeInTheDocument();

    // if we click on the radio button, it should be checked

    await userEvent.click(radioButton);

    expect(radioButton).toBeChecked();
    const wrapper2 = getByTestId('tx-minimum-threshold');
    expect(wrapper2).toBeInTheDocument();

    // const input1 = within(wrapper2).getByRole('number');
    // expect(input1).toBeInTheDocument();

    // userEvent.type(input1, '1000');
    // expect(input1).toHaveValue('1000');
    // const input2 = within(wrapper2).getByRole('combobox');
    // expect(input2).toBeInTheDocument();
    // userEvent.selectOptions(input2, 'USD');
    // expect(input2).toHaveValue('USD');

    // if we click again, it should be unchecked
    // userEvent.click(radioButton);
    // expect(radioButton).not.toBeChecked();
  });
});
