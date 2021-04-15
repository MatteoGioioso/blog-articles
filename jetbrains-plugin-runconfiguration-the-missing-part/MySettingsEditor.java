import com.github.matteogioioso.rbijetbrainsplugin.toolWindow.runConfigurations.MyRunConfiguration;
import com.intellij.execution.configuration.EnvironmentVariablesTextFieldWithBrowseButton;
import com.intellij.openapi.options.ConfigurationException;
import com.intellij.openapi.options.SettingsEditor;
import com.intellij.openapi.ui.ComponentWithBrowseButton;
import com.intellij.openapi.ui.LabeledComponent;
import com.intellij.openapi.ui.TextFieldWithBrowseButton;
import com.intellij.ui.components.JBLabel;
import org.jetbrains.annotations.NotNull;

import javax.swing.*;

public class MySettingsEditor extends SettingsEditor<MyRunConfiguration> {
    private LabeledComponent<ComponentWithBrowseButton> MainComponent;
    private EnvironmentVariablesTextFieldWithBrowseButton envVarsField;
    private JTextField commandField;

    @Override
    protected void resetEditorFrom(@NotNull MyRunConfiguration s) {
        envVarsField.setEnvs(s.getEnvironmentalVariables());

    }

    @Override
    protected void applyEditorTo(@NotNull MyRunConfiguration s) throws ConfigurationException {
        s.setEnvironmentalVariables(envVarsField.getEnvs());
        s.setMyCommand(commandField.getText());
    }

    @Override
    protected @NotNull JComponent createEditor() {
        JPanel jPanel = new JPanel();
        JBLabel envVarLabel = new JBLabel("Environmental variables:");
        this.envVarsField = new EnvironmentVariablesTextFieldWithBrowseButton();

        JBLabel commandLabel = new JBLabel("Input a command:");
        this.commandField = new JTextField();

        jPanel.add(envVarLabel);
        jPanel.add(envVarsField);
        jPanel.add(commandLabel);
        jPanel.add(commandField);
        return jPanel;
    }

    private void createUIComponents() {
        MainComponent = new LabeledComponent<>();
        MainComponent.setComponent(new TextFieldWithBrowseButton());
    }
}
