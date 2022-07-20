import { dynamicActivate } from 'utils/i18nLoaderHelper';
import { render, screen } from 'utils/test-utils';
import PasswordStrength from '.';

describe('<PasswordStrength />', () => {
  beforeEach(() => {
    dynamicActivate('en');
  });

  describe('contains9Characters', () => {
    it('should be green color', () => {
      render(<PasswordStrength data="xxxx xxxxx" />);
      expect(screen.getByTestId('contains9Characters__icon')).toHaveStyle({ color: 'green' });
    });

    it('should be gray color', () => {
      render(<PasswordStrength data="xxxx xx" />);
      expect(screen.getByTestId('contains9Characters__icon')).toHaveStyle({ color: 'gray.900' });
    });
  });

  describe('containsOneLowerCase', () => {
    it('should be green color', () => {
      render(<PasswordStrength data="nxx xxx o" />);
      expect(screen.getByTestId('containsOneLowerCase__icon')).toHaveStyle({ color: 'green' });
    });

    it('should be gray color', () => {
      render(<PasswordStrength data="000" />);
      expect(screen.getByTestId('containsOneLowerCase__icon')).toHaveStyle({ color: 'gray.200' });
    });
  });

  describe('containsOneUpperCase', () => {
    it('should be green color', () => {
      render(<PasswordStrength data="nNx xxx o" />);
      expect(screen.getByTestId('containsOneUpperCase__icon')).toHaveStyle({ color: 'green' });
    });

    it('should be gray color', () => {
      render(<PasswordStrength data="nx xxx o" />);
      expect(screen.getByTestId('containsOneUpperCase__icon')).toHaveStyle({ color: 'gray.900' });
    });
  });

  describe('containsOneNumber', () => {
    it('should be green color', () => {
      render(<PasswordStrength data="nNx xxx 12" />);
      expect(screen.getByTestId('containsOneNumber__icon')).toHaveStyle({ color: 'green' });
    });

    it('should be gray color', () => {
      render(<PasswordStrength data="nx xxx o" />);
      expect(screen.getByTestId('containsOneNumber__icon')).toHaveStyle({ color: 'gray.900' });
    });
  });

  describe('containsOneSpecialChar', () => {
    it('should be green color', () => {
      render(<PasswordStrength data="nNx xxx 12_!" />);
      expect(screen.getByTestId('containsOneSpecialChar__icon')).toHaveStyle({ color: 'green' });
    });

    it('should be gray color', () => {
      render(<PasswordStrength data="nx xxx o" />);
      expect(screen.getByTestId('containsOneSpecialChar__icon')).toHaveStyle({ color: 'gray.900' });
    });
  });
});
