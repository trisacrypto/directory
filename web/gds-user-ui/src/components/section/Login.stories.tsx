import React from "react";
import { Story } from "@storybook/react";
import Login from "./Login";

interface LoginProps {}

export default {
  title: "Components/Login",
  component: Login,
};

export const standard: Story<LoginProps> = ({ ...props }) => (
  <Login {...props} />
);

standard.bind({});
