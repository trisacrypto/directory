import React from "react";
import { Story } from "@storybook/react";
import AuthEmailConfirmation from './AuthEmailConfirmation';

interface AuthEmailConfirmationProps {
}

export default {
  title: "Components/AuthEmailConfirmation",
  component: AuthEmailConfirmation,
};

export const standard: Story<AuthEmailConfirmationProps> = ({ ...props }) => (
  <AuthEmailConfirmation {...props} />
);

standard.bind({
});
