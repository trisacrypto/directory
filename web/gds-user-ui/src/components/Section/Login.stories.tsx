import React from 'react';
import { Story } from '@storybook/react';
import Login from './Login';
interface LoginProps {
  handleSignWithSocial: (event: React.FormEvent, type: string) => void;
  handleSignWithEmail: (data: any) => void;
  isLoading?: boolean;
  isError?: any;
}

export default {
  title: 'Components/Login',
  component: Login
};

export const standard: Story<LoginProps> = ({ ...props }) => <Login {...props} />;

standard.bind({});
