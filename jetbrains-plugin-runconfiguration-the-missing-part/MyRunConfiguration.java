package com.github.matteogioioso.rbijetbrainsplugin.toolWindow.runConfigurations;

import com.github.matteogioioso.rbijetbrainsplugin.toolWindow.views.RunDebuggerSettingsEditor;
import com.intellij.execution.ExecutionException;
import com.intellij.execution.Executor;
import com.intellij.execution.configurations.*;
import com.intellij.execution.process.OSProcessHandler;
import com.intellij.execution.process.ProcessHandler;
import com.intellij.execution.process.ProcessHandlerFactory;
import com.intellij.execution.process.ProcessTerminatedListener;
import com.intellij.execution.runners.ExecutionEnvironment;
import com.intellij.openapi.options.SettingsEditor;
import com.intellij.openapi.project.Project;
import org.jetbrains.annotations.NotNull;
import org.jetbrains.annotations.Nullable;

import java.util.ArrayList;
import java.util.Map;

public class MyRunConfiguration extends RunConfigurationBase<MyRunConfigurationOptions> {

    protected MyRunConfiguration(@NotNull Project project, @Nullable ConfigurationFactory factory, @Nullable String name) {
        super(project, factory, name);
    }

    @Override
    public @NotNull SettingsEditor<? extends RunConfiguration> getConfigurationEditor() {
        return new RunDebuggerSettingsEditor();
    }

    @Override
    public @Nullable RunProfileState getState(@NotNull Executor executor, @NotNull ExecutionEnvironment environment) throws ExecutionException {
        return new CommandLineState(environment) {
            /**
             * This method will run echo and output whatever you input in the command field
             */
            @Override
            protected @NotNull ProcessHandler startProcess() throws ExecutionException {
                ArrayList<String> strings = new ArrayList<>();
                strings.add("echo");
                strings.add(getOptions().getMyCommand());
                GeneralCommandLine commandLine = new GeneralCommandLine(strings);

                OSProcessHandler processHandler = ProcessHandlerFactory.getInstance().createColoredProcessHandler(commandLine);
                ProcessTerminatedListener.attach(processHandler);
                return processHandler;
            }
        };
    }

    @NotNull
    @Override
    protected MyRunConfigurationOptions getOptions() {
        return (MyRunConfigurationOptions) super.getOptions();
    }

    public Map<String, String> getEnvironmentalVariables() {
        return getOptions().getEnvironmentalVariables();
    }

    public void setEnvironmentalVariables(Map<String, String> env) {
        getOptions().setEnvironmentalVariables(env);
    }

    public String getMyCommand() {
        return getOptions().getMyCommand();
    }

    public void setMyCommand(String command) {
        getOptions().setMyCommand(command);
    }
}
