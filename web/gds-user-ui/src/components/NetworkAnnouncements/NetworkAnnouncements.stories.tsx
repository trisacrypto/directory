import { Meta, Story } from "@storybook/react";
import NetworkAnnouncements from ".";

interface NetworkAnnouncementsProps {
  message: string;
}
export default {
  title: "components/NetworkAnnouncements",
  component: NetworkAnnouncements,
} as Meta;

const Template: Story = (args: NetworkAnnouncementsProps) => (
  <NetworkAnnouncements {...args} />
);

export const Default = Template.bind({});
Default.args = {};
