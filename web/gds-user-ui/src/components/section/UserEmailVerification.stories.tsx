import React from "react";
import { Story } from "@storybook/react";
import UserEmailVerification from './UserEmailVerification';

interface UserEmailVerificationProps {
}

export default {
  title: "Components/UserEmailVerification",
  component: UserEmailVerification,
};

export const standard: Story<UserEmailVerificationProps> = ({ ...props }) => (
  <UserEmailVerification {...props} />
);

standard.bind({
});
