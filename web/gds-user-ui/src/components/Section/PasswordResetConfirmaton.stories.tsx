import React from "react";
import { Story } from "@storybook/react";
import PasswordResetConfirmation from './PasswordResetConfirmation';

interface PasswordResetConfirmationProps {
  email: String
}

export default {
  title: "Components/PasswordResetConfirmation",
  component: PasswordResetConfirmation,
};

export const standard: Story<PasswordResetConfirmationProps> = ({ ...props }) => (
  <PasswordResetConfirmation {...props} />
);

standard.bind({
  email: 'cletus@100kode.io'
});
