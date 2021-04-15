import java.util.HashMap;
import java.util.Map;
import com.intellij.execution.ExecutionManager;
import com.intellij.execution.RunManager;
import com.intellij.execution.RunnerAndConfigurationSettings;
import com.intellij.execution.executors.DefaultRunExecutor;
import com.intellij.execution.runners.ExecutionEnvironmentBuilder;
import com.intellij.openapi.project.Project;

public class RunConfigurationPromptFunction {
    /**
     *
     * @param project Jetbrains Project
     */
    private promptRunConfiguration(Project project) {
        return errorDomain -> {
            Map<String, String> environmentalVariables = new HashMap<>();
            environmentalVariables.put("MY_ENV", "hello");

            MyRunConfigurationType myRunConfigurationType = new MyRunConfigurationType();
            RunManager instance = RunManager.getInstance(project);
            RunnerAndConfigurationSettings configurationTemplate = instance
                    .getConfigurationTemplate(myRunConfigurationType.getConfigurationFactories()[0]);

            // The most important part is to cast the configuration
            // from the template to your custom configuration
            MyRunConfiguration configuration = (MyRunConfiguration) configurationTemplate.getConfiguration();
            configuration.setEnvironmentalVariables(environmentalVariables);

            ExecutionEnvironmentBuilder builder = ExecutionEnvironmentBuilder
                    .createOrNull(DefaultRunExecutor.getRunExecutorInstance(), nodejsTemplate);

            if (builder != null) {
                ExecutionManager.getInstance(project).restartRunProfile(builder.build());
            }
        };
    }
}
