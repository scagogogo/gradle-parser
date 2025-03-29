// 一个示例Kotlin DSL格式的Gradle构建文件

plugins {
    kotlin("jvm") version "1.7.10"
    kotlin("plugin.spring") version "1.7.10"
    id("org.springframework.boot") version "2.7.0"
    id("io.spring.dependency-management") version "1.0.11.RELEASE"
}

group = "com.example"
version = "0.1.0-SNAPSHOT"
description = "示例Kotlin Spring Boot项目"

java {
    sourceCompatibility = JavaVersion.VERSION_11
    targetCompatibility = JavaVersion.VERSION_11
}

repositories {
    mavenCentral()
    google()
    maven { url = uri("https://jitpack.io") }
    maven {
        url = uri("https://maven.aliyun.com/repository/public")
        // 示例仓库凭证
        credentials {
            username = "user"
            password = "password"
        }
    }
}

dependencies {
    // Spring Boot依赖
    implementation("org.springframework.boot:spring-boot-starter-web")
    implementation("org.springframework.boot:spring-boot-starter-data-jpa")
    
    // Kotlin支持
    implementation("org.jetbrains.kotlin:kotlin-reflect")
    implementation("org.jetbrains.kotlin:kotlin-stdlib-jdk8")
    
    // 数据库驱动
    implementation("mysql:mysql-connector-java:8.0.29")
    
    // 工具库
    implementation("org.apache.commons:commons-lang3:3.12.0")
    implementation("com.google.guava:guava:31.1-jre")
    
    // 测试依赖
    testImplementation("org.springframework.boot:spring-boot-starter-test")
    testImplementation("org.junit.jupiter:junit-jupiter-api:5.8.2")
    testRuntimeOnly("org.junit.jupiter:junit-jupiter-engine:5.8.2")
    
    // 示例项目内依赖
    implementation(project(":common"))
}

// 自定义任务
tasks.register("showDependencies") {
    group = "custom"
    description = "显示所有依赖"
    
    doLast {
        println("项目依赖列表:")
        configurations.implementation.get().allDependencies.forEach { dependency ->
            println("- ${dependency.group}:${dependency.name}:${dependency.version}")
        }
    }
}

// 示例扩展配置
springBoot {
    mainClass.set("com.example.ApplicationKt")
}

tasks.test {
    useJUnitPlatform()
} 