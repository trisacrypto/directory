import React from "react";
import { Story } from "@storybook/react";
import UserEmailConfirmation from './UserEmailConfirmation';

interface UserEmailConfirmationProps {
}

export default {
  title: "Components/UserEmailConfirmation",
  component: UserEmailConfirmation,
};

export const standard: Story<UserEmailConfirmationProps> = ({ ...props }) => (
  <UserEmailConfirmation {...props} />
);

standard.bind({
});
