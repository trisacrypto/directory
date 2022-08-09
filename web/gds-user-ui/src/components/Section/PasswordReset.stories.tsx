import React from 'react';
import { Story } from '@storybook/react';
import PasswordReset from './PasswordReset';

interface PasswordResetProps {
  handleSubmit: (data: any) => void;
  isLoading: boolean;
  isError?: any;
}

export default {
  title: 'Components/PasswordReset',
  component: PasswordReset
};

export const standard: Story = ({ ...props }) => <PasswordReset {...props} />;

standard.bind({});
