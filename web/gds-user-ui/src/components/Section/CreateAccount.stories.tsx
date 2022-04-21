import React from 'react';
import { Story } from '@storybook/react';
import CreateAccount from './CreateAccount';
interface CreateAccountProps {
  handleSignUp: (evt: any, type: any) => void;
}
export default {
  title: 'Components/CreateAccount',
  component: CreateAccount
};

export const standard: Story<CreateAccountProps> = ({ ...props }) => <CreateAccount {...props} />;

standard.bind({});
