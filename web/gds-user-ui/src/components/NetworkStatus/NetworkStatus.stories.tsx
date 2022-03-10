import { Meta, Story } from "@storybook/react";
import NetworkStatus from ".";

export default {
  title: "components/NetworkStatus",
  component: NetworkStatus,
} as Meta;

const Template: Story = (args) => <NetworkStatus {...args} />;

export const Default = Template.bind({});
Default.args = {
  isOnline: false,
};
