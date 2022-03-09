import { Meta, Story } from "@storybook/react";
import NetworkStatus from ".";

interface NetworkStatusProps {
  isOnline: boolean;
}
export default {
  title: "components/NetworkStatus",
  component: NetworkStatus,
} as Meta;

const Template: Story = (args: NetworkStatusProps) => (
  <NetworkStatus {...args} />
);

export const Default = Template.bind({});
Default.args = {
  isOnline: false,
};
