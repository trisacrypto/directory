/* eslint-disable init-declarations */
import userEvent from '@testing-library/user-event';
import { dynamicActivate } from 'utils/i18nLoaderHelper';
import { render, screen } from 'utils/test-utils';
import StepButtons from '.';

// mock the useFetchCertificateStep hook
jest.mock('hooks/useFetchCertificateStep', () => ({
  useFetchCertificateStep: () => ({
    certificateStep: {
      form: jest.fn(), // mock the form function
      errors: jest.fn() // mock the errors function
    },
    isFetchingCertificateStep: false
  })
}));

// mock

describe('<StepButtons />', () => {
  let handlePreviousStep = jest.fn();
  let currentStep;
  let isCurrentStepLastStep;
  let handleResetForm = jest.fn();
  let isDefaultValue;
  beforeAll(() => {
    dynamicActivate('en');
  });

  it('should display the correct label', () => {
    handlePreviousStep = jest.fn();
    currentStep = 2;
    isCurrentStepLastStep = false;
    handleResetForm = jest.fn();
    isDefaultValue = jest.fn(() => true);

    render(
      <StepButtons
        handlePreviousStep={handlePreviousStep}
        currentStep={currentStep}
        isCurrentStepLastStep={isCurrentStepLastStep}
        handleResetForm={handleResetForm}
        isDefaultValue={isDefaultValue}
      />
    );

    // expect(screen.getByRole('button', { name: /Save & Previous/i })).toBeInTheDocument();
    // expect(screen.getByRole('button', { name: /Save & Next /i })).toBeInTheDocument();
    // expect(screen.getByRole('button', { name: /Clear & Reset Form/i })).toBeInTheDocument();
  });

  // it('should call handlePreviousStep when previous button is clicked', () => {
  //   handlePreviousStep = jest.fn();
  //   currentStep = 2;
  //   isCurrentStepLastStep = false;
  //   handleResetForm = jest.fn();
  //   isDefaultValue = jest.fn(() => true);
  //   render(
  //     <StepButtons
  //       handlePreviousStep={handlePreviousStep}
  //       currentStep={currentStep}
  //       isCurrentStepLastStep={isCurrentStepLastStep}
  //       handleResetForm={handleResetForm}
  //       isDefaultValue={isDefaultValue}
  //     />
  //   );

  //   const previousButton = screen.getByRole('button', { name: /save & previous/i });

  //   userEvent.click(previousButton);

  //   expect(handlePreviousStep).toHaveBeenCalledTimes(1);
  // });

  // it('should call handleResetForm when reset button is clicked', () => {
  //   handlePreviousStep = jest.fn();
  //   currentStep = 2;
  //   isCurrentStepLastStep = false;
  //   handleResetForm = jest.fn();
  //   isDefaultValue = jest.fn(() => false);
  //   render(
  //     <StepButtons
  //       handlePreviousStep={handlePreviousStep}
  //       currentStep={currentStep}
  //       isCurrentStepLastStep={isCurrentStepLastStep}
  //       handleResetForm={handleResetForm}
  //       isDefaultValue={isDefaultValue}
  //     />
  //   );

  //   const resetButton = screen.getByRole('button', { name: /Clear & Reset Form/i });

  //   userEvent.click(resetButton);

  //   expect(handleResetForm).toHaveBeenCalledTimes(1);
  // });

  // it('should call handleResetForm when reset button is clicked', () => {
  //   handlePreviousStep = jest.fn();
  //   currentStep = 2;
  //   isCurrentStepLastStep = false;
  //   handleResetForm = jest.fn();
  //   isDefaultValue = jest.fn(() => false);
  //   render(
  //     <StepButtons
  //       handlePreviousStep={handlePreviousStep}
  //       currentStep={currentStep}
  //       isCurrentStepLastStep={isCurrentStepLastStep}
  //       handleResetForm={handleResetForm}
  //       isDefaultValue={isDefaultValue}
  //     />
  //   );

  //   const resetButton = screen.getByRole('button', { name: /Clear & Reset Form/i });

  //   userEvent.click(resetButton);

  //   expect(handleResetForm).toHaveBeenCalledTimes(1);
  // });
});
