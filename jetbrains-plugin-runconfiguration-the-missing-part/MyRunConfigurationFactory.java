package com.github.matteogioioso.rbijetbrainsplugin.toolWindow.runConfigurations;

import com.intellij.execution.configurations.ConfigurationFactory;
import com.intellij.execution.configurations.RunConfiguration;
import com.intellij.openapi.components.BaseState;
import com.intellij.openapi.project.Project;
import org.jetbrains.annotations.NotNull;
import org.jetbrains.annotations.Nullable;

public class MyRunConfigurationFactory extends ConfigurationFactory {
    private static final String FACTORY_NAME = "Run debugger configuration factory";

    public MyRunConfigurationFactory(MyRunConfigurationType type) {
        super(type);
    }

    @Override
    public @NotNull RunConfiguration createTemplateConfiguration(@NotNull Project project) {
        return new MyRunConfiguration(project, this, "Run");
    }

    @Override
    public @NotNull String getId() {
        return "my-runconfiguration";
    }

    @Override
    public @NotNull String getName() {
        return "MyRunConfiguration";
    }

    @Nullable
    @Override
    public Class<? extends BaseState> getOptionsClass() {
        return MyRunConfigurationOptions.class;
    }
}
