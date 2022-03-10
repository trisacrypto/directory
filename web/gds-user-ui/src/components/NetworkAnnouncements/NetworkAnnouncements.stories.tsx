import { Meta, Story } from "@storybook/react";
import NetworkAnnouncements from ".";

export default {
  title: "components/NetworkAnnouncements",
  component: NetworkAnnouncements,
} as Meta;

const Template: Story = (args) => <NetworkAnnouncements {...args} />;

export const Default = Template.bind({});
Default.args = {};
