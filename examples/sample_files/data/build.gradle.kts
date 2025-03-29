// data模块的Gradle构建文件 (Kotlin DSL)

plugins {
    id("java-library")
    kotlin("jvm") version "1.7.10"
}

group = "com.example.data"
version = rootProject.extra["appVersion"] as String
description = "数据处理模块"

java {
    sourceCompatibility = JavaVersion.VERSION_11
    targetCompatibility = JavaVersion.VERSION_11
}

repositories {
    mavenCentral()
}

dependencies {
    // 项目模块依赖
    implementation(project(":common"))
    
    // 数据库依赖
    api("org.hibernate:hibernate-core:5.6.5.Final")
    api("org.springframework.data:spring-data-jpa:2.7.0")
    implementation("com.h2database:h2:2.1.214")
    
    // Kotlin标准库
    implementation("org.jetbrains.kotlin:kotlin-stdlib-jdk8")
    
    // 测试依赖
    testImplementation("org.junit.jupiter:junit-jupiter-api:5.8.2")
    testRuntimeOnly("org.junit.jupiter:junit-jupiter-engine:5.8.2")
}

tasks.test {
    useJUnitPlatform()
} 