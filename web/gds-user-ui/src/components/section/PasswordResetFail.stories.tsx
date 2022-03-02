import React from "react";
import { Story } from "@storybook/react";
import PasswordResetFail from './PasswordResetFail';

interface PasswordResetFailProps {
  email: string;
}

export default {
  title: "Components/PasswordResetFail",
  component: PasswordResetFail,
};

export const standard: Story<PasswordResetFailProps> = ({ ...props }) => (
  <PasswordResetFail {...props} />
);

standard.bind({
  email: 'cletus@100kode.io'
});
