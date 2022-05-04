import React from "react";
import { Story } from "@storybook/react";
import PasswordReset from "./PasswordReset";

interface PasswordResetProps {}

export default {
  title: "Components/PasswordReset",
  component: PasswordReset,
};

export const standard: Story<PasswordResetProps> = ({ ...props }) => (
  <PasswordReset {...props} />
);

standard.bind({});
