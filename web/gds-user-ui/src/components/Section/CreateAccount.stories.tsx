import React from 'react';
import { Story } from '@storybook/react';
import CreateAccount from './CreateAccount';
interface CreateAccountProps {
  handleSocialAuth: (event: React.FormEvent, type: string) => void;
  handleSignUpWithEmail: (data: any) => void;
}
export default {
  title: 'Components/CreateAccount',
  component: CreateAccount
};

export const standard: Story<CreateAccountProps> = ({ ...props }) => <CreateAccount {...props} />;

standard.bind({});
