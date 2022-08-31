import userEvent from '@testing-library/user-event';
import { dynamicActivate } from 'utils/i18nLoaderHelper';
import { render, screen } from 'utils/test-utils';
import ReviewSubmit from '.';

describe('ReviewSubmit', () => {
  beforeEach(() => {
    dynamicActivate('en');
    localStorage.clear();
  });
  it('should call handleSubmitRegister with testnet if we click on testnet submitting button', () => {
    const handleSubmitRegister = jest.fn();
    render(<ReviewSubmit onSubmitHandler={handleSubmitRegister} />);

    const testnetSubmittingButtonEl = screen.getByTestId('testnet-submit-btn');

    userEvent.click(testnetSubmittingButtonEl);

    expect(handleSubmitRegister).toHaveBeenCalledTimes(1);
    expect(handleSubmitRegister.mock.calls[0].length).toBe(2);
    expect(handleSubmitRegister.mock.calls[0]).toContain('testnet');
  });

  it('should call handleSubmitRegister with mainnet if we click on mainnet submitting button', () => {
    const handleSubmitRegister = jest.fn();
    render(<ReviewSubmit onSubmitHandler={handleSubmitRegister} />);

    const mainnetSubmittingButtonEl = screen.getByTestId('mainnet-submit-btn');

    userEvent.click(mainnetSubmittingButtonEl);

    expect(handleSubmitRegister).toHaveBeenCalledTimes(1);
    expect(handleSubmitRegister.mock.calls[0].length).toBe(2);
    expect(handleSubmitRegister.mock.calls[0]).toContain('mainnet');
  });

  // 'TODO: refactor this test'

  // it('should disable testnet submitting button when testnet data are sent', () => {
  //   localStorage.setItem('isTestNetSent', JSON.stringify(true));
  //   const handleSubmitRegister = jest.fn();
  //   render(<ReviewSubmit onSubmitHandler={handleSubmitRegister} />);

  //   const testnetSubmittingButtonEl = screen.getByTestId('testnet-submit-btn');

  //   expect(localStorage.getItem).toHaveBeenCalled();
  //   expect(testnetSubmittingButtonEl).toBeDisabled();
  // });

  // it('should enable testnet submitting button when testnet data are not sent', () => {
  //   localStorage.setItem('isTestNetSent', JSON.stringify(false));
  //   const handleSubmitRegister = jest.fn();
  //   render(<ReviewSubmit onSubmitHandler={handleSubmitRegister} />);

  //   const testnetSubmittingButtonEl = screen.getByTestId('testnet-submit-btn');

  //   expect(localStorage.getItem).toHaveBeenCalled();
  //   expect(testnetSubmittingButtonEl).toBeEnabled();
  // });

  // it('should disable testnet submitting button when testnet data are sent', () => {
  //   localStorage.setItem('isMainNetSent', JSON.stringify(true));
  //   const handleSubmitRegister = jest.fn();
  //   render(<ReviewSubmit onSubmitHandler={handleSubmitRegister} />);

  //   const mainnetSubmittingButtonEl = screen.getByTestId('mainnet-submit-btn');

  //   expect(localStorage.getItem).toHaveBeenCalled();
  //   expect(mainnetSubmittingButtonEl).toBeDisabled();
  // });

  // it('should disable testnet submitting button when testnet data are sent', () => {
  //   localStorage.setItem('isMainNetSent', JSON.stringify(false));
  //   const handleSubmitRegister = jest.fn();
  //   render(<ReviewSubmit onSubmitHandler={handleSubmitRegister} />);

  //   const mainnetSubmittingButtonEl = screen.getByTestId('mainnet-submit-btn');

  //   expect(localStorage.getItem).toHaveBeenCalled();
  //   expect(mainnetSubmittingButtonEl).toBeEnabled();
  // });
});
